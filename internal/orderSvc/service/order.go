package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/needon1997/theshop-svc/internal/orderSvc/model"
	"github.com/needon1997/theshop-svc/internal/orderSvc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

var GoodsSvcConn *grpc.ClientConn
var InventorySvcConn *grpc.ClientConn

type OrderService struct {
	proto.UnimplementedOrderServer
}

func panicIfErr(desc string, err error) {
	if err != nil {
		var wrapperErr error
		errDesc := fmt.Sprintf("[%s error]: %s", desc, err.Error())
		switch err {
		case gorm.ErrRecordNotFound:
			wrapperErr = status.Error(codes.NotFound, errDesc)
		default:
			wrapperErr = status.Error(codes.Internal, errDesc)
		}
		panic(wrapperErr)
	}
}
func (OrderService) CartItemList(ctx context.Context, in *proto.UserInfo) (rsp *proto.CartItemListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CartItemList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cartItemList := make([]model.ShoppingCart, 0)
	panicIfErr("Match User", model.DB.Where("user_id = ?", in.Id).Find(&cartItemList).Error)
	rsp = &proto.CartItemListResponse{}
	rsp.Total = int32(len(cartItemList))
	for i := 0; i < int(rsp.Total); i++ {
		cartItem := &proto.ShopCartInfoResponse{
			Id:      int32(cartItemList[i].ID),
			UserId:  int32(cartItemList[i].UserId),
			GoodsId: int32(cartItemList[i].GoodsId),
			Nums:    int32(cartItemList[i].Num),
			Checked: cartItemList[i].Checked,
		}
		rsp.Data = append(rsp.Data, cartItem)
	}
	return
}
func (OrderService) CreateCartItem(ctx context.Context, in *proto.CartItemRequest) (rsp *proto.ShopCartInfoResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateCartItem]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cartItem := model.ShoppingCart{}
	err1 := model.DB.Where("user_id = ?", in.UserId).Where("goods_id = ?", in.GoodsId).First(&cartItem).Error
	if err1 != nil && err1 != gorm.ErrRecordNotFound {
		panicIfErr("Find CartItem", err1)
	}
	if err1 == gorm.ErrRecordNotFound {
		cartItem.Checked = in.Checked
		cartItem.UserId = uint(in.UserId)
		cartItem.GoodsId = uint(in.GoodsId)
		cartItem.Num = uint(in.Nums)
		panicIfErr("Create CartItem", model.DB.Create(&cartItem).Error)
	} else {
		cartItem.Checked = in.Checked
		cartItem.Num += uint(in.Nums)
		panicIfErr("Merge CartItem", model.DB.Save(&cartItem).Error)
	}
	rsp = &proto.ShopCartInfoResponse{
		Id:      int32(cartItem.ID),
		UserId:  int32(cartItem.UserId),
		GoodsId: int32(cartItem.GoodsId),
		Nums:    int32(cartItem.Num),
		Checked: cartItem.Checked,
	}
	return
}
func (OrderService) UpdateCartItem(ctx context.Context, in *proto.CartItemRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateCartItem]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cartItem := model.ShoppingCart{}
	panicIfErr("Find CartItem", model.DB.Where("user_id = ?", in.UserId).Where("goods_id = ?", in.GoodsId).First(&cartItem).Error)
	cartItem.Checked = in.Checked
	if in.Nums > 0 {
		cartItem.Num = uint(in.Nums)
	}
	panicIfErr("Save CartItem", model.DB.Save(&cartItem).Error)
	rsp = &empty.Empty{}
	return
}
func (OrderService) DeleteCartItem(ctx context.Context, in *proto.CartItemRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteCartItem]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cartItem := model.ShoppingCart{}
	panicIfErr("Find CartItem", model.DB.Where("user_id = ?", in.UserId).Where("goods_id = ?", in.GoodsId).First(&cartItem).Error)
	panicIfErr("Delete Item", model.DB.Delete(&cartItem).Error)
	rsp = &empty.Empty{}
	return
}
func (OrderService) CreateOrder(ctx context.Context, in *proto.OrderRequest) (rsp *proto.OrderInfoResponse, err error) {
	tx := model.DB.Begin()
	defer func() {
		tx.Rollback()
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateOrder]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cartItem := make([]model.ShoppingCart, 0)
	panicIfErr("Find Cart Item", tx.Where("user_id = ?", in.UserId).Where("checked = ?", "true").Find(cartItem).Error)
	goodsClient := proto.NewGoodsClient(GoodsSvcConn)
	goodsId := make([]int32, 0)
	for i := 0; i < len(cartItem); i++ {
		goodsId = append(goodsId, int32(cartItem[i].GoodsId))
	}
	goods, err1 := goodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsId})
	if err1 != nil {
		panic(err1)
	}
	var orderTotal float32 = 0.0
	for i := 0; i < len(goods.Data); i++ {
		orderTotal += goods.Data[i].ShopPrice * float32(cartItem[i].Num)
	}
	inventoryClient := proto.NewInventoryClient(InventorySvcConn)
	sellInfo := &proto.SellInfo{}
	for i := 0; i < len(cartItem); i++ {
		sellInfo.GoodsInfo = append(sellInfo.GoodsInfo, &proto.GoodsInvInfo{
			GoodsId: int32(cartItem[i].GoodsId),
			Num:     int32(cartItem[i].Num),
		})
	}
	_, err1 = inventoryClient.Sell(context.Background(), sellInfo)
	if err1 != nil {
		panic(err1)
	}
	order := model.OrderInfo{
		UserId:      uint(in.UserId),
		OrderSn:     uuid.New().String(),
		OrderAmount: orderTotal,
		PayAt:       time.Time{},
		Address:     in.Address,
		BuyerName:   in.Name,
		BuyerMobile: in.Mobile,
		Note:        in.Note,
	}
	panicIfErr("Create Order", model.DB.Create(&order).Error)
	orderGoodsList := make([]model.OrderGoods, 0)
	for i := 0; i < len(goodsId); i++ {
		orderGoodsList = append(orderGoodsList, model.OrderGoods{
			OrderId:       order.ID,
			GoodsId:       uint(goodsId[i]),
			GoodsName:     goods.Data[i].Name,
			GoodsImageURL: goods.Data[i].GoodsFrontImage,
			GoodsPrice:    goods.Data[i].ShopPrice,
			Num:           cartItem[i].Num,
		})
	}
	panicIfErr("Create Order Goods", model.DB.Create(&orderGoodsList).Error)
	panicIfErr("Find Cart Item", tx.Where("user_id = ?", in.UserId).Where("checked = ?", "true").Delete(&model.ShoppingCart{}).Error)
	panicIfErr("Commit Change", tx.Commit().Error)
	rsp = &proto.OrderInfoResponse{
		Id:      int32(order.ID),
		UserId:  int32(order.UserId),
		OrderSn: order.OrderSn,
		Status:  order.Status,
		Note:    order.Note,
		Total:   order.OrderAmount,
		Address: order.Address,
		Name:    order.BuyerName,
		Mobile:  order.BuyerMobile,
		AddTime: order.CreatedAt.String(),
	}
	return
}
func (OrderService) OrderList(ctx context.Context, in *proto.OrderFilterRequest) (rsp *proto.OrderListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[OrderList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	orderList := make([]model.OrderInfo, 0)
	db := model.DB
	if in.UserId != 0 {
		db = model.DB.Where("user_id = ?", in.UserId)
	}
	panicIfErr("List Order", db.Find(&orderList).Error)
	rsp = &proto.OrderListResponse{}
	rsp.Total = int32(len(orderList))
	psize := 10
	pn := 1
	if in.PagePerNums != 0 {
		psize = int(in.PagePerNums)
	}
	if in.Pages != 0 {
		pn = int(in.Pages)
	}
	offset := (pn - 1) * psize
	for i := offset; i < offset+psize; i++ {
		order := &proto.OrderInfoResponse{
			Id:      int32(orderList[i].ID),
			UserId:  int32(orderList[i].UserId),
			OrderSn: orderList[i].OrderSn,
			Status:  orderList[i].Status,
			Note:    orderList[i].Note,
			Total:   orderList[i].OrderAmount,
			Address: orderList[i].Address,
			Name:    orderList[i].BuyerName,
			Mobile:  orderList[i].BuyerMobile,
			AddTime: orderList[i].CreatedAt.String(),
		}
		rsp.Data = append(rsp.Data, order)
	}
	return
}
func (OrderService) OrderDetail(ctx context.Context, in *proto.OrderRequest) (rsp *proto.OrderInfoDetailResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[OrderDetail]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	order := &model.OrderInfo{}
	if in.UserId != 0 {
		panicIfErr("Find Order", model.DB.Where("user_id = ?", in.UserId).Where("id = ?", in.Id).First(&order).Error)
	} else {
		panicIfErr("Find Order", model.DB.First(&order, in.Id).Error)
	}
	orderGoods := make([]model.OrderGoods, 0)
	panicIfErr("Find OrderGoods", model.DB.Where("order_id = ?", order.ID).Find(&orderGoods).Error)
	rsp = &proto.OrderInfoDetailResponse{}
	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:      int32(order.ID),
		UserId:  int32(order.UserId),
		OrderSn: order.OrderSn,
		Status:  order.Status,
		Note:    order.Note,
		Total:   order.OrderAmount,
		Address: order.Address,
		Name:    order.BuyerName,
		Mobile:  order.BuyerMobile,
		AddTime: order.CreatedAt.String(),
	}

	total := len(orderGoods)
	for i := 0; i < total; i++ {
		rsp.Data = append(rsp.Data, &proto.OrderItemResponse{
			Id:         int32(orderGoods[i].ID),
			OrderId:    int32(orderGoods[i].OrderId),
			GoodsId:    int32(orderGoods[i].GoodsId),
			GoodsName:  orderGoods[i].GoodsName,
			GoodsImage: orderGoods[i].GoodsImageURL,
			GoodsPrice: orderGoods[i].GoodsPrice,
			Nums:       int32(orderGoods[i].Num),
		})
	}
	return
}
func (OrderService) UpdateOrderStatus(ctx context.Context, in *proto.OrderStatus) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateOrderStatus]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	order := &model.OrderInfo{}
	panicIfErr("Find Order", model.DB.Where("order_sn = ?", in.OrderSn).First(&order).Error)
	order.Status = in.Status
	panicIfErr("Update Status", model.DB.Save(&order).Error)
	rsp = &empty.Empty{}
	return
}
