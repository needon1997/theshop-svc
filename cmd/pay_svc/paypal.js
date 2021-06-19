var paypal = require('paypal-rest-sdk');
var grpc = require('@grpc/grpc-js');
paypal.configure({
    mode: 'sandbox', // Sandbox or live
    client_id: '',
    client_secret: ''
});
module.exports.create = function create_payment(call, callback) {
    console.log("1");
    var order_sn = call.request.order_sn;
    var payReq = JSON.stringify({
        intent: 'sale',
        payer: {
            payment_method: 'paypal'
        },
        redirect_urls: {
            return_url: `http://localhost:10081/o/v1/payment/execute/${order_sn}?x-token=${call.request.token}`,
            cancel_url: `http://localhost:10081/o/v1/order/${order_sn}?x-token=${call.request.token}`,
        },
        transactions: [{
            amount: {
                total: call.request.total/100,
                currency: call.request.currency,
            },
        }]
    });
    paypal.payment.create(payReq, function (error, payment) {
        var links = {};
        if (error) {
            return callback({code: grpc.status.INTERNAL, message: error});
        } else {
            // Capture HATEOAS links
            payment.links.forEach(function (linkObj) {
                links[linkObj.rel] = {
                    href: linkObj.href,
                    method: linkObj.method
                };
            })

            // If the redirect URL is present, redirect the customer to that URL
            if (links.hasOwnProperty('approval_url')) {
                return callback(null, {accept_url: links.approval_url.href})
            } else {
                return callback({code: grpc.status.INTERNAL, message: "unknown payment error"});
            }
        }
    });
}
module.exports.execute = function execute_payment(call, callback) {
    var paymentId = call.request.payment_id;
    var payerId = {payer_id: call.request.payer_id};

    paypal.payment.execute(paymentId, payerId, function (error, payment) {
        if (error) {
            console.error(JSON.stringify(error));
        } else {
            if (payment.state == 'approved') {
                return callback(null, {msg: 'payment completed successfully', status: payment.state})
            } else {
                return callback({code: grpc.status.INTERNAL, message: "payment not successful"})
            }
        }
    });
}
