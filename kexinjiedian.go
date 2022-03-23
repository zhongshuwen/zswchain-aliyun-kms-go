package main

import (
	"context"
	"fmt"

	zsw "github.com/zhongshuwen/zswchain-go"
	zswattr "github.com/zhongshuwen/zswchain-go/zswattr"
	"github.com/zhongshuwen/zswchain-go/zswitems"
	"github.com/zhongshuwen/zswchain-go/zswperms"
)

type ItemBalanceTableRow struct {
	ItemId           uint64 `json:"item_id"`
	Status           uint32 `json:"status"`
	Balance          uint64 `json:"balance"`
	BalanceInCustody uint64 `json:"balance_in_custody"`
	BalanceFrozen    uint64 `json:"balance_frozen"`
}

func QueryUserCangpin(ctx context.Context, api *zsw.API, account zsw.AccountName) (out *[]ItemBalanceTableRow, errOut error) {
	var rowReq = zsw.GetTableRowsRequest{
		Code:       "zsw.items",
		Scope:      string(account),
		Table:      "itembalances",
		LowerBound: "", //use this to paginate with the last result's id
		Limit:      10, //results to fetch
		JSON:       true,
	}
	var resp, err = api.GetTableRows(ctx, rowReq)
	if err != nil {
		errOut = err
		return
	}
	var x []ItemBalanceTableRow
	resp.JSONToStructs(&x)
	return &x, nil

}
func GetCreateExampleCollection(authorizer zsw.AccountName, kexinJiedian zsw.AccountName, collectionZswId string) []*zsw.Action {
	collectionSchemaName := zsw.Name("v1collection")

	itemMode := zsw.ITEM_CONFIG_ALLOW_MUTABLE_DATA |
		zsw.ITEM_CONFIG_ALLOW_NOTIFY |
		zsw.ITEM_CONFIG_TRANSFERABLE

	return []*zsw.Action{
		/*
			zswitems.NewMakeSchema(
				authorizer,
				kexinJiedian,
				collectionSchemaName,
				[]zsw.FieldDef{
					{Name: "name", Type: "string"},
					{Name: "description", Type: "string"},
					{Name: "logo_url", Type: "string"},
					{Name: "website", Type: "string"},
					{Name: "banner_url", Type: "string"},
					{Name: "icon_url", Type: "string"},
				},
			),*/
		zswitems.NewMakeCollection(
			authorizer,
			UuidToUint128OrQuit(collectionZswId),
			UuidToUint128OrQuit(collectionZswId).GetTypeACode(),
			0,
			kexinJiedian,
			kexinJiedian,
			uint32(itemMode),
			500,
			500,
			kexinJiedian,
			10000000,
			10000000,
			0,
			collectionSchemaName,
			[]zsw.AccountName{kexinJiedian},
			[]zsw.AccountName{},
			[]zsw.AccountName{kexinJiedian},
			zswattr.AttributeMap{
				"name":       zswattr.ToZSWAttribute("My Collection"),
				"logo_url":   zswattr.ToZSWAttribute("https://zhongshuwen.com/wp-content/uploads/2021/08/zswlogo@4x.png"),
				"website":    zswattr.ToZSWAttribute("https://zhongshuwen.com"),
				"banner_url": zswattr.ToZSWAttribute("https://zhongshuwen.com/wp-content/uploads/2021/08/zswlogo@4x.png"),
			},
			fmt.Sprintf("https://dummydata.testnet.chao7.cn/dd/collections/%d.json", UuidToUint128OrQuit(collectionZswId).Get40BitId()),
		),
	}
}
func GetActionsCreateExampleCollectionItemFlow1155(authorizer zsw.AccountName, kexinJiedian zsw.AccountName, collectionZswId string, itemTemplateZswId string, itemZswId string) []*zsw.Action {
	itemSchemaName := zsw.Name(RandomLowercaseStringAZ(12))

	itemMode :=
		zsw.ITEM_CONFIG_ALLOW_NOTIFY |
			zsw.ITEM_CONFIG_TRANSFERABLE

	return []*zsw.Action{

		zswitems.NewMakeSchema(
			authorizer,
			kexinJiedian,
			itemSchemaName,
			[]zsw.FieldDef{
				{Name: "name", Type: "string"},
				{Name: "image_url", Type: "string"},
				{Name: "rarity", Type: "uint32"},
			},
		),

		zswitems.NewMakeItemTemplate(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(itemTemplateZswId),
			UuidToUint128OrQuit(itemTemplateZswId).GetTypeACode(),
			UuidToUint128OrQuit(collectionZswId).GetTypeACode(),
			0,
			itemSchemaName,
			zswattr.AttributeMap{},
			"",
		),

		zswitems.NewMakeItem(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(itemZswId).Get40BitId(),
			UuidToUint128OrQuit(itemZswId),
			uint32(itemMode),
			UuidToUint128OrQuit(itemTemplateZswId).GetTypeACode(),
			9000000,
			itemSchemaName,
			zswattr.AttributeMap{
				"name":      zswattr.ToZSWAttribute("很酷的数字藏品"),
				"image_url": zswattr.ToZSWAttribute("https://cangpin.test.chao7.cn/f/images/shanghai.png"),
			},
			zswattr.AttributeMap{},
		),
	}
}
func GetActionsCreateExampleCollectionItemFlow721(authorizer zsw.AccountName, kexinJiedian zsw.AccountName, collectionZswId string, itemTemplateZswId string, itemZswId string) []*zsw.Action {
	itemSchemaName := zsw.Name(RandomLowercaseStringAZ(12))
	collectionSchemaName := zsw.Name("simpleimage1")

	itemMode := zsw.ITEM_CONFIG_ALLOW_MUTABLE_DATA |
		zsw.ITEM_CONFIG_ALLOW_NOTIFY |
		zsw.ITEM_CONFIG_TRANSFERABLE

	return []*zsw.Action{
		zswitems.NewMakeSchema(
			authorizer,
			kexinJiedian,
			itemSchemaName,
			[]zsw.FieldDef{
				{Name: "name", Type: "string"},
				{Name: "logo_url", Type: "string"},
				{Name: "website", Type: "string"},
				{Name: "banner_url", Type: "string"},
			},
		),
		zswitems.NewMakeCollection(
			authorizer,
			UuidToUint128OrQuit(collectionZswId),
			UuidToUint128OrQuit(collectionZswId).GetTypeACode(),
			0,
			kexinJiedian,
			kexinJiedian,
			uint32(itemMode),
			500,
			500,
			kexinJiedian,
			10000000,
			10000000,
			0,
			collectionSchemaName,
			[]zsw.AccountName{kexinJiedian},
			[]zsw.AccountName{},
			[]zsw.AccountName{kexinJiedian},
			zswattr.AttributeMap{
				"name":       zswattr.ToZSWAttribute("My Collection"),
				"logo_url":   zswattr.ToZSWAttribute("https://zhongshuwen.com/wp-content/uploads/2021/08/zswlogo@4x.png"),
				"website":    zswattr.ToZSWAttribute("https://zhongshuwen.com"),
				"banner_url": zswattr.ToZSWAttribute("https://zhongshuwen.com/wp-content/uploads/2021/08/zswlogo@4x.png"),
			},
			fmt.Sprintf("https://dummydata.testnet.chao7.cn/dd/collections/%d.json", UuidToUint128OrQuit(collectionZswId).Get40BitId()),
		),
		zswitems.NewMakeSchema(
			authorizer,
			kexinJiedian,
			itemSchemaName,
			[]zsw.FieldDef{
				{Name: "name", Type: "string"},
				{Name: "image_url", Type: "string"},
				{Name: "rarity", Type: "uint32"},
				{Name: "xp", Type: "uint64"},
			},
		),

		zswitems.NewMakeItemTemplate(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(itemTemplateZswId),
			UuidToUint128OrQuit(itemTemplateZswId).GetTypeACode(),
			UuidToUint128OrQuit(collectionZswId).GetTypeACode(),
			0,
			itemSchemaName,
			zswattr.AttributeMap{},
			"",
		),

		zswitems.NewMakeItem(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(itemZswId).Get40BitId(),
			UuidToUint128OrQuit(itemZswId),
			uint32(itemMode),
			UuidToUint128OrQuit(itemTemplateZswId).GetTypeACode(),
			1,
			itemSchemaName,
			zswattr.AttributeMap{
				"name":      zswattr.ToZSWAttribute("很酷的数字藏品"),
				"image_url": zswattr.ToZSWAttribute("https://cangpin.test.chao7.cn/f/images/shanghai.png"),
			},
			zswattr.AttributeMap{
				"xp": zswattr.ToZSWAttribute(uint64(100)),
			},
		),
	}
}

func GetActionsSetupKexinJiedianPermissions(authorizer zsw.AccountName, kexinJiedian zsw.AccountName, kexinJiedianZswId string) []*zsw.Action {
	return []*zsw.Action{

		// 给可信节点Minting权限
		zswperms.NewSetZswPerms(
			authorizer,    //中数文的内容审核管理账号
			kexinJiedian,  //可信节点联盟链账号
			"zsw.prmcore", //core permissions scope
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TX)| // 此权限赋予客户用户授权交易的权力
					uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TRANSFER_ITEM), //允许可信节点赋予C2C基本数字藏品转移
			),
		),
		// 给可信节点自愿监护权限
		zswitems.NewMakeCustodian(
			authorizer,                             //中数文的内容审核管理账号
			kexinJiedian,                           //平台生成的walletName
			UuidToUint128OrQuit(kexinJiedianZswId), //中数文平台的“userId”（登录借口获取的）
			zsw.NewUint128FromUint64(0),            //现在0，没有用
			zsw.NewUint128FromUint64(
				uint64(zsw.CUSTODIAN_PERMS_ENABLED)| //开通Custodian功能
					uint64(zsw.CUSTODIAN_PERMS_TX_TO_SELF_CUSTODIAN)| //可以authorize用户在自己的
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_NULL_CUSTODIAN)| //can send from self custodianship to another custodian
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_ZSW_CUSTODIAN), //can send from self custodianship to a non-custodial null custodian
			),
			0, //0是征程
			0, //其他的可信节点用户要使用你的平台的时候，数字藏品要冻多久（秒）
			[]zsw.AccountName{
				kexinJiedian, //为了查看历史方便，可以设置logevent账号，未来也可以加handler
			},
		),
		zswitems.NewMakeIssuer(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(kexinJiedianZswId),
			zsw.NewUint128FromUint64(0),
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_ITEM)| //允许基本minting的功能
					uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_TO_NULL_CUSTODIAN), //可以mint到需要用户公钥权限的custodian
			),
			0, //0==正常
		),
		zswitems.NewMakeRoyaltyUser( //登记谁是版税接受者
			authorizer, //需要中数文签名
			kexinJiedian,
			UuidToUint128OrQuit(kexinJiedianZswId),
			zsw.NewUint128FromUint64(0),
			0,
		),
	}
}

func GetActionsCreateDummyCollection(authorizer zsw.AccountName, kexinJiedian zsw.AccountName, kexinJiedianZswId string) []*zsw.Action {
	return []*zsw.Action{

		// 给可信节点Minting权限
		zswperms.NewSetZswPerms(
			authorizer,    //中数文的内容审核管理账号
			kexinJiedian,  //可信节点联盟链账号
			"zsw.prmcore", //core permissions scope
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TX)| // 此权限赋予客户用户授权交易的权力
					uint64(zsw.ZSW_CORE_PERMS_CONFIRM_AUTHORIZE_USER_TRANSFER_ITEM), //允许可信节点赋予C2C基本数字藏品转移
			),
		),
		// 给可信节点自愿监护权限
		zswitems.NewMakeCustodian(
			authorizer,                             //中数文的内容审核管理账号
			kexinJiedian,                           //平台生成的walletName
			UuidToUint128OrQuit(kexinJiedianZswId), //中数文平台的“userId”（登录借口获取的）
			zsw.NewUint128FromUint64(0),            //现在0，没有用
			zsw.NewUint128FromUint64(
				uint64(zsw.CUSTODIAN_PERMS_ENABLED)| //开通Custodian功能
					uint64(zsw.CUSTODIAN_PERMS_TX_TO_SELF_CUSTODIAN)| //可以authorize用户在自己的
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_NULL_CUSTODIAN)| //can send from self custodianship to another custodian
					uint64(zsw.CUSTODIAN_PERMS_SEND_TO_ZSW_CUSTODIAN), //can send from self custodianship to a non-custodial null custodian
			),
			0, //0是征程
			0, //其他的可信节点用户要使用你的平台的时候，数字藏品要冻多久（秒）
			[]zsw.AccountName{
				kexinJiedian, //为了查看历史方便，可以设置logevent账号，未来也可以加handler
			},
		),
		zswitems.NewMakeIssuer(
			authorizer,
			kexinJiedian,
			UuidToUint128OrQuit(kexinJiedianZswId),
			zsw.NewUint128FromUint64(0),
			zsw.NewUint128FromUint64(
				uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_ITEM)| //允许基本minting的功能
					uint64(zsw.ZSW_ITEMS_PERMS_AUTHORIZE_MINT_TO_NULL_CUSTODIAN), //可以mint到需要用户公钥权限的custodian
			),
			0, //0==正常
		),
		zswitems.NewMakeRoyaltyUser( //登记谁是版税接受者
			authorizer, //需要中数文签名
			kexinJiedian,
			UuidToUint128OrQuit(kexinJiedianZswId),
			zsw.NewUint128FromUint64(0),
			0,
		),
	}
}
