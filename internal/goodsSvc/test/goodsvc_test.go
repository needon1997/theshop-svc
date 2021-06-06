package test

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"google.golang.org/grpc"
	"testing"
)

var client proto.GoodsClient

func init() {
	conn, _ := grpc.Dial("127.0.0.1:10084", grpc.WithInsecure())
	client = proto.NewGoodsClient(conn)
}
func TestGetGoodListFilter(t *testing.T) {
	list, err := client.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin: 50,
		PriceMax: 500,
		IsHot:    true,
		IsNew:    true,
	})
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(list.Data); i++ {
		fmt.Println(list.Data[i].IsHot, list.Data[i].IsNew)
	}
}
func TestGetGoodsListById(t *testing.T) {
	list, err := client.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: []int32{2, 3, 4, 5, 6, 7}})
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(list.Data); i++ {
		fmt.Println(list.Data[i])
	}
}

func TestGetAllCatList(t *testing.T) {
	list, err := client.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list.JsonData)
}
