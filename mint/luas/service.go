package luas

type IService interface {
	Sha256()
	//	hash 函数,sha256
	//	签名验证 函数,sigToAddress
	//	字符串转整形函数
	//	整形转字符串函数
	//	转大写
	//	转小写
}

type Service struct {
}

func (s *Service) Sha256(data ... []byte) {

}
