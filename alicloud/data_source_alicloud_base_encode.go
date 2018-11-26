package alicloud

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudBaseEncode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBaseEncodeRead,

		Schema: map[string]*schema.Schema{
			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"encode_data": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"padding": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudBaseEncodeRead(d *schema.ResourceData, meta interface{}) error {
	var keyword = "dxc28DXCdxc39DXC"
	var iv = "dxc14DXCdxc69DXC"

	data := d.Get("data").(string)
	if keywordInput, ok := d.GetOk("keyword"); ok {
		if len(keywordInput.(string)) != 16 {
			return fmt.Errorf("Base encode got an error: %s", "The length of keyword must be 16.")
		}
		keyword = keywordInput.(string)
	}

	encodeDataString, padding, err := encrypt([]byte(data), []byte(keyword), []byte(iv))
	if err != nil {
		return fmt.Errorf("Base encode got an error: %#v", "The length of keyword must be 16.", err)
	}

	d.Set("encode_data", encodeDataString)
	d.Set("padding", padding)
	d.SetId(encodeDataString)

	return nil
}

func encrypt(origData []byte, key []byte, iv []byte) (string, int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", 0, err
	}
	blockSize := block.BlockSize()
	origData, padding := pkcs5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), padding, nil
}

func decrypt(crypted string, key []byte, iv []byte) (string, error) {
	decodeData, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(origData, decodeData)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) ([]byte, int) {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...), padding
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)

	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
