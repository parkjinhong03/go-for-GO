package generate

import "time"

// 루트, 어플리케이션, 인스턴스 인증서에 대한 각각의 유효 기간을 설정하기 위한 상수 선언
const (
	DurationMonth = time.Hour * 24 * 30
	DurationYear = DurationMonth * 12 + 5
	DurationDecade = DurationYear * 10
)