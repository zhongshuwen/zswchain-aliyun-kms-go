package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zhongshuwen/gmsm/sm2"
	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandomLowercaseStringAZ(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func checkSigCore2(signature string, digest []byte) error {
	sigNew, err := ecc.NewSignature(signature)
	if err != nil {
		fmt.Printf("checkSigCore2: checkSig error %w\n", err)
		return err
	}
	pubKey1 := ecc.MustNewPublicKeyFromData(sigNew.Content[0:33])
	pubKey2 := ecc.MustNewPublicKey("PUB_GM_6VmANYA7h8VwU4dbEeC6dTbGzYxRwukyW7BMz6Zsc93NUbwPRA")
	fmt.Printf("key 1: %s, key 2: %s\n", pubKey1.String(), pubKey2.String())
	res1 := sigNew.Verify(digest, pubKey1)
	res2 := sigNew.Verify(digest, pubKey2)
	fmt.Printf("res1: %t\nres2: %t\n", res1, res2)
	decomp := sm2.Decompress(sigNew.Content[0:33])
	len := uint32(sigNew.Content[34]) + 33
	result := decomp.VerifyDigest(digest, sigNew.Content[33:len])
	if !result {
		return fmt.Errorf("checkSigCore2: verify digest failed")
	}
	return nil
}
func main() {

	//zsw.EnableDebugLogging(zsw.NewLogger(false))

	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "allkms",
			Aliases: []string{},
			Usage:   "run kms tests on user",
			Action: func(c *cli.Context) error {
				_, err := RunDebugScenarioC(c.Context, "kmsalikmsali", zsw.AccountName(RandomLowercaseStringAZ(12)))
				return err
			},
		},
		{
			Name:    "keys",
			Aliases: []string{"k"},
			Subcommands: []*cli.Command{{
				Name:        "generate",
				Aliases:     []string{"gen"},
				Usage:       "生成 N 个国密PM2私钥+公钥",
				UsageText:   "zswkmsdemo keys generate 100",
				Description: "生成 N 个国密PM2私钥+公钥",
				ArgsUsage:   "<n = 1>",
				Category:    "",
				Action: func(c *cli.Context) error {
					countStr := c.Args().First()
					generateCount := 1
					if len(countStr) != 0 {
						candidateCount, err := strconv.Atoi(c.Args().First())
						if err != nil {
							return fmt.Errorf("invalid setting for n %w", err)
						} else if candidateCount <= 0 {
							return fmt.Errorf("invalid setting for n, must be larger than 0")
						} else {
							generateCount = candidateCount
						}
					}
					for i := 0; i < generateCount; i++ {
						privateKey, err := ecc.NewRandomPrivateKey()
						if err != nil {
							return fmt.Errorf("error generating private key %w", err)
						}
						fmt.Printf("================================================================\n密钥: %s\n公钥: %s\n", privateKey.String(), privateKey.PublicKey().String())
					}
					return nil
				},
				Subcommands:            []*cli.Command{},
				Flags:                  []cli.Flag{},
				SkipFlagParsing:        false,
				HideHelp:               false,
				HideHelpCommand:        false,
				Hidden:                 false,
				UseShortOptionHandling: false,
				HelpName:               "",
				CustomHelpTemplate:     "",
			}, {
				Name:        "convert",
				Aliases:     []string{"conv"},
				Usage:       "convert key format x to key format y",
				UsageText:   "zswkmsdemo keys convert",
				Description: "",
				Category:    "",
				Action: func(c *cli.Context) error {
					_, err := ConvertFile([]byte{}, c.String("from"), c.String("to"), c.String("i"), c.String("o"))

					return err
				},
				Subcommands: []*cli.Command{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "from-format",
						Aliases:  []string{"from"},
						Usage:    "<pubpem | pubzswkey | privzswkey | rawprivkey | comppubkey>",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to-format",
						Aliases:  []string{"to"},
						Usage:    "<pubzswkey | privzswkey>",
						Required: true,
					},
					&cli.StringFlag{
						Name:      "input-file",
						Aliases:   []string{"i"},
						Usage:     "Input File",
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:    "output-file",
						Aliases: []string{"o"},
						Usage:   "Input File",
					},
				},
				SkipFlagParsing:        false,
				HideHelp:               false,
				HideHelpCommand:        false,
				Hidden:                 false,
				UseShortOptionHandling: false,
				HelpName:               "",
				CustomHelpTemplate:     "",
			}},
			Flags:                  []cli.Flag{},
			SkipFlagParsing:        false,
			HideHelp:               false,
			HideHelpCommand:        false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
			CustomHelpTemplate:     "",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
