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

func (m *MongoStore) Search(name string) []Kitten {
	s := m.session.Clone()
	defer s.Close()

	var result []Kitten
	// 아래 코드는 kittenserver 데이터베이스에서 kittens 컬렉션을 검색하여 반환한다.
	c := s.DB("kittenserver").C("kittens")
	// Find 메서드를 통해 원하는 값 검색 후 All 메서드를 이용해 결과값을 매개변수에 대입한다.
	err := c.Find(Kitten{Name: name}).All(&result)
	if err != nil {
		return nil
	}

	return result
}

// kittenserver 데이터베이스 있는 모든 kitten 콜렉션 정보를 삭제하는 메서드
func (m *MongoStore) DeleteAllKittens() {
	s := m.session.Clone()
	defer s.Close()

	_ = s.DB("kittenserver").C("kittens").DropCollection()
}

// kittenserver 데이터베이스에 매개변수로 넘겨 받은 객체를 kittens 콜렉션에 저장하는 메서드
func (m *MongoStore) InsertKittens(kittens []Kitten) {
	s := m.session.Clone()
	defer s.Close()

	_ = s.DB("kittenserver").C("kittens").Insert(kittens)
}