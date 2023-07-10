package webvpn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Webvpn struct {
	Enable bool   `yaml:"enable"` // 使能 Webvpn
	Server string `yaml:"Server"` // Webvpn 内网服务器地址
	Host   string `yaml:"Host"`   // Webvpn 外网接收域名
	Port   int    `yaml:"Port"`   // Webvpn 端口
	Tls    bool   `yaml:"Tls"`    // 访问 Webvpn 是否使用TLS加密（https）
	AesKey string `yaml:"AesKey"` // 普通链接转换 Webvpn 链接的AesKey
	AesIv  string `yaml:"AesIv"`  // 普通链接转换 Webvpn 链接的AesIv
	Cookie string `yaml:"Cookie"` // 访问 Webvpn 是否使用TLS加密（https）
}

var UrlMap = make(map[string]string)
var Cfg Webvpn

type PersonJson struct {
	CanVisitConnection bool     `json:"canVisitConnection"`
	CanVisitProtocol   []string `json:"canVisitProtocol"`
	InitPassword       bool     `json:"initPassword"`
	SelfAccess         bool     `json:"selfAccess"`
	ShowFAQ            bool     `json:"showFAQ"`
	UserType           string   `json:"userType"`
	Username           string   `json:"username"`
	WeekPassword       bool     `json:"weekPassword"`
	WrdvpnIV           string   `json:"wrdvpnIV"`
	WrdvpnKey          string   `json:"wrdvpnKey"`
}

func LoadWebvpn() error {
	// log.Debug("我来证明你解析了WebvpnCfg", Cfg)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	baseUrl := "https://" + Cfg.Server + "/user/info"
	err, body := getWebVpnInfo(baseUrl, client)
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(time.Minute * 3)
			// 第二次开始, 每3分钟 get一次，防止cookie失效
			err, body := getWebVpnInfo(baseUrl, client)
			var p PersonJson
			err = json.Unmarshal(body, &p)
			if err != nil {
				log.Error("webvpn返回了非json数据，可能是cookie过期了，请更新！")
			}
		}
	}()

	// 第一次, 解析json数据
	var p PersonJson
	err = json.Unmarshal(body, &p)
	if err != nil {
		return errors.New("webvpn返回了非json数据，可能是cookie过期了，请更新！")
	}
	if len(p.Username) > 8 && len(p.Username) < 15 {
		log.Info("获取webvpn数据成功，cookie的学号是：", p.Username)
		Cfg.AesKey = p.WrdvpnKey
		Cfg.AesIv = p.WrdvpnIV
		log.Info("更新webvpn的AES密钥：Key:", Cfg.AesKey, " ,IV:", Cfg.AesIv)
	} else {
		return errors.New("无法从webvpn获取数据，可能是cookie无效了，请更换")
	}
	log.Debug("获取webvpn数据", string(body))
	return nil
}

func getWebVpnInfo(baseUrl string, client *http.Client) (error, []byte) {
	req, err := http.NewRequest("GET", baseUrl, nil)
	req.Header.Set("Host", Cfg.Host)
	req.Header.Add("Cookie", Cfg.Cookie)
	if err != nil {
		return err, nil
	}
	response, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err, nil
	}
	return err, body
}

func UrlEncoding(RawUrl string) (string, error) {
	if _, ok := UrlMap[RawUrl]; ok {
		// 存在
		// log.Debug("url在map里面")
		return UrlMap[RawUrl], nil
	} else {
		// log.Debug("url不在map里面")
		// 不存在
		byteData := []byte(Cfg.AesIv)
		// 将 byte 装换为 16进制的字符串
		hexStringData := hex.EncodeToString(byteData)
		var encrypted, err = AesEncryptCFB([]byte(RawUrl), []byte(Cfg.AesKey), []byte(Cfg.AesIv))
		if err != nil {
			return "", err
		}
		hexStringData = hexStringData + encrypted
		// byte 转 16进制 的结果
		log.Debug("加密的url是", hexStringData)
		UrlMap[RawUrl] = hexStringData
		return hexStringData, nil
	}
}

func AesEncryptCFB(origData []byte, key []byte, iv []byte) (string, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// The IV needs to be unique, but not secure. Therefore, it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(origData))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], origData)
	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	// return fmt.Sprintf("%x", ciphertext), nil
	return strings.TrimLeft(hex.EncodeToString(ciphertext), "0"), nil
}
