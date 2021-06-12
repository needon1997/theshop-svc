const paypal = require('./paypal');
var PROTO_PATH = __dirname + "/proto/payment.proto";
var grpc = require('@grpc/grpc-js');
var protoLoader = require('@grpc/proto-loader');
var health = require('./health');
var axios = require('axios')
const statusMap = {"": proto.grpc.health.v1.HealthCheckResponse.ServingStatus.SERVING};
let healthImpl = new health.Implementation(statusMap);
// Suggested options for similarity to existing grpc.load behavior
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
var protoDescriptor = grpc.loadPackageDefinition(packageDefinition);

function getServer() {
    var server = new grpc.Server();
    server.addService(health.service, healthImpl);
    server.addService(protoDescriptor.Payment.service, {
        CreatePayment: paypal.create,
        ExecutePayment: paypal.execute
    });
    return server;
}

var routeServer = getServer();
routeServer.bindAsync('0.0.0.0:10083', grpc.ServerCredentials.createInsecure(), () => {
    axios.put('http://127.0.0.1:8500/v1/agent/service/register', {
            Name: "pay-service",
            Address: '127.0.0.1',
            Port: 10083,
            ID: 'pay-service',
            Check: {
                GRPC: '127.0.0.1:10083',
                Interval: '20s'
            }
        },
        {
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(res => {
        console.log(`statusCode: ${res.status}`)
    }).catch(error => {
        console.error(error)
    })
    routeServer.start();
});
