package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
	"github.com/zhongshuwen/zswchain-go/zswitems"
	"github.com/zhongshuwen/zswchain-tencent-kms-go/kmswallet"
)

func RunDebugScenarioA(ctx context.Context, api *zsw.API, creator zsw.AccountName, newKexinJiedian zsw.AccountName, newKexinJiedianPublicKey string, newKexinJiedianZswId string) (error, string) {
	actions := []*zsw.Action{}
	actions = append(actions, GetActionsCreateUserWithResources(
		creator,
		newKexinJiedian,
		newKexinJiedianPublicKey,
		1000000,
		zsw.NewZSWAsset(10000),
		zsw.NewZSWAsset(2000),
	)...)
	actions = append(actions, GetActionsSetupKexinJiedianPermissions(
		creator,
		newKexinJiedian,
		newKexinJiedianZswId,
	)...)

	return nil, runTxBasic(context.Background(), api, actions)
}

func RunDebugScenarioB(ctx context.Context, api *zsw.API, authorizer zsw.AccountName, kexinJiedian zsw.AccountName, kexinJiedianPublicKey string, kexinJiedianZswId string, recipientUser zsw.AccountName) (error, string) {
	collectionZswId := uuid.New().String()
	itemTemplateZswId := uuid.New().String()
	itemZswId := uuid.New().String()

	actions := GetActionsCreateExampleCollectionItemFlow1155(
		authorizer,
		kexinJiedian,
		collectionZswId,
		itemTemplateZswId,
		itemZswId,
	)
	mintActions := []*zsw.Action{
		zswitems.NewItemMint(
			kexinJiedian,
			recipientUser,
			kexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{1},
			"An item for you!",
		),
	}

	actions = append(actions, mintActions...)

	return nil, runTxBasic(context.Background(), api, actions)
}
func getRandUUID() string {
	b := make([]byte, 4)

	rand.Read(b)
	return fmt.Sprintf("00000000-0000-0000-0000-0000%s", hex.EncodeToString(b))
	//rand.Intn(100)

}

func RunDebugScenarioC(ctx context.Context, authorizer zsw.AccountName, newKexinJiedian zsw.AccountName) (string, error) {

	api := zsw.New("https://node3.tn1.chao7.cn")
	//api.Debug = true
	client, err := kmswallet.GetKMSClient(
		os.Getenv("ALIYUN_KMS_AK_ID"),
		os.Getenv("ALIYUN_KMS_AK_SECRET"),
		"ap-hangzhou",
		"kms.tencentcloudapi.com",
	)
	if err != nil {
		return "", err
	}
	keyBag := kmswallet.NewAliyunKMSKeyBag(client)
	zswKey, err := keyBag.AddKMSKeyById(os.Getenv("ALIYUN_KMS_KEY_ID"), os.Getenv("ALIYUN_KMS_KEY_VERSION_ID"))
	fmt.Printf("added key %s from KMS\n", zswKey)
	if err != nil {
		return "", err
	}
	/*
		NoError(
			keyBag.ImportPrivateKeyFromEnv(context.Background(), "ZSW_CONTENT_REVIEW_PRIVATE_KEY"),
			"missing ZSW_CONTENT_REVIEW_PRIVATE_KEY",
		)
	*/
	kexinJiedianZswId := uuid.New().String()
	collectionZswId := uuid.New().String()
	itemTemplateZswId := uuid.New().String()
	itemZswId := uuid.New().String()

	kxjdPrivateKey, err := ecc.NewPrivateKey("PVT_GM_E23jvM1z35D4UxfYTmWLS9ButJwXJ13zHuZwvUjpxwqEVQLPX")
	if err != nil {
		return "", err
	}
	userAPrivateKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	userBPrivateKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	keyBag.Append(kxjdPrivateKey)

	userAName := zsw.AccountName(fmt.Sprintf("usra1%s", RandomLowercaseStringAZ(7)))
	userBName := zsw.AccountName(fmt.Sprintf("usrb1%s", RandomLowercaseStringAZ(7)))

	fmt.Printf("-- 新可信节点 --\n账号：%s\n公钥：%s\n密钥：%s\n--------------------------------------------------\n", newKexinJiedian, kxjdPrivateKey.PublicKey().String(), kxjdPrivateKey.String())
	fmt.Printf("-- 用户A --\n账号：%s\n公钥：%s\n密钥：%s\n--------------------------------------------------\n", userAName, userAPrivateKey.PublicKey().String(), userAPrivateKey.String())
	fmt.Printf("-- 用户B --\n账号：%s\n公钥：%s\n密钥：%s\n--------------------------------------------------\n", userBName, userBPrivateKey.PublicKey().String(), userBPrivateKey.String())

	api.SetSigner(keyBag)
	actions := []*zsw.Action{}
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		newKexinJiedian,
		kxjdPrivateKey.PublicKey().String(),
		1000000,
		zsw.NewZSWAsset(10000),
		zsw.NewZSWAsset(2000),
	)...)
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		userAName,
		userAPrivateKey.PublicKey().String(),
		3000,
		zsw.NewZSWAsset(0),
		zsw.NewZSWAsset(0),
	)...)
	actions = append(actions, GetActionsCreateUserWithResources(
		authorizer,
		userBName,
		userBPrivateKey.PublicKey().String(),
		3000,
		zsw.NewZSWAsset(0),
		zsw.NewZSWAsset(0),
	)...)
	actions = append(actions, GetActionsSetupKexinJiedianPermissions(
		authorizer,
		newKexinJiedian,
		kexinJiedianZswId,
	)...)
	runTxBasic(context.Background(), api, actions)
	time.Sleep(time.Second * 2)
	actions = GetCreateExampleCollection(
		authorizer,
		newKexinJiedian,
		collectionZswId,
	)
	runTxBasic(context.Background(), api, actions)
	time.Sleep(time.Second * 2)

	actions = GetActionsCreateExampleCollectionItemFlow1155(
		authorizer,
		newKexinJiedian,
		collectionZswId,
		itemTemplateZswId,
		itemZswId,
	)
	runTxBasic(context.Background(), api, actions)
	time.Sleep(time.Second * 2)
	mintActions := []*zsw.Action{

		zswitems.NewItemMint(
			newKexinJiedian,
			userAName,
			newKexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{100},
			"An item for you!",
		),
		zswitems.NewItemMint(
			newKexinJiedian,
			userBName,
			newKexinJiedian,
			0,
			[]uint64{UuidToUint128OrQuit(itemZswId).Get40BitId()},
			[]uint64{200},
			"An item for you 2!",
		),
	}

	//actions = append(actions, mintActions...)

	return runTxBasic(context.Background(), api, mintActions), nil

}
