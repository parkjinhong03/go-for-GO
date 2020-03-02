package data

import "gopkg.in/mgo.v2"

// 모의 객체가 아닌, 실제 데이터 저장소인 MongoDB를 DataStore 인터페이스로 사용하기 위한 구조체 정의
type MongoStore struct {
	session *mgo.Session
}

// 매개변수로 받은 IP로 MongoDB 세션을 연결시키고, 그 연결시킨 세션을 필드에 넣어 MongoStore 구조체를 반환하는 메서드.
func NewMongoStore(connection string) (*MongoStore, error) {
	// mgo.Dial 함수를 이용하여 MongoDB 서버에 연결시킬 수 있다.
	session, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}

	return &MongoStore{session: session}, nil
}
