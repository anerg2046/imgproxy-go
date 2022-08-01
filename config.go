package imgproxygo

type Config struct {
	BaseUrl       string //图片服务域名
	Key           string //签名KEY
	Salt          string //签名SALT
	SignatureSize int    //签名长度，需与IMGPROXY_SIGNATURE_SIZE一致
	Encode        bool   //是否对目标图片地址进行base64编码
}
