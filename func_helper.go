/* Copyright (C) 2015-2019 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015-2019년 UnHa Kim (unha.kim@kuh.pe.kr)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package xing

// #include "./types_c.h"
import "C"

import (
	"fmt"
	"github.com/ghts/lib"
	"github.com/mitchellh/go-ps"
	"unsafe"

	"strings"
	"time"
)

func F전일() time.Time {
	return 전일.G값()
}

func F당일() time.Time {
	return 당일.G값()
}

func F최근_영업일_모음() []time.Time {
	return lib.F슬라이스_복사(최근_영업일_모음, nil).([]time.Time)
}

func F2전일_시각(포맷 string, 값 interface{}) (time.Time, error) {
	if strings.Contains(포맷, "2") {
		return time.Time{}, lib.New에러("포맷에 이미 날짜가 포함되어 있습니다. %v", 포맷)
	}

	시각, 에러 := lib.F2포맷된_시각(포맷, 값)
	if 에러 != nil {
		return time.Time{}, 에러
	}

	전일 := F전일()

	전일_시각 := time.Date(전일.Year(), 전일.Month(), 전일.Day(),
		시각.Hour(), 시각.Minute(), 시각.Second(), 0, 시각.Location())

	return 전일_시각, nil
}

func F2전일_시각_단순형(포맷 string, 값 interface{}) time.Time {
	return lib.F확인(F2전일_시각(포맷, 값)).(time.Time)
}

func F2당일_시각(포맷 string, 값 interface{}) (time.Time, error) {
	if strings.Contains(포맷, "2") {
		return time.Time{}, lib.New에러("포맷에 이미 날짜가 포함되어 있습니다. %v", 포맷)
	}

	시각, 에러 := lib.F2포맷된_시각(포맷, 값)
	if 에러 != nil {
		return time.Time{}, 에러
	}

	당일 := F당일()

	당일_시각 := time.Date(당일.Year(), 당일.Month(), 당일.Day(),
		시각.Hour(), 시각.Minute(), 시각.Second(), 0, 시각.Location())

	return 당일_시각, nil
}

func F2당일_시각_단순형(포맷 string, 값 interface{}) time.Time {
	return lib.F확인(F2당일_시각(포맷, 값)).(time.Time)
}

func xing_C32_실행_중() (프로세스ID int) {
	defer lib.S예외처리{M함수: func() { 프로세스ID = -1 }}.S실행()

	프로세스_모음, 에러 := ps.Processes()
	lib.F확인(에러)

	for _, 프로세스 := range 프로세스_모음 {
		if 실행화일명 := 프로세스.Executable(); strings.HasSuffix(xing_C32_경로, 실행화일명) {
			return 프로세스.Pid()
		}
	}

	return -1
}

//func xing_COM32_실행_중() (프로세스ID int) {
//	defer lib.S예외처리{M함수: func() { 프로세스ID = -1 }}.S실행()
//
//	프로세스_모음, 에러 := ps.Processes()
//	lib.F에러체크(에러)
//
//	for _, 프로세스 := range 프로세스_모음 {
//		if 실행화일명 := 프로세스.Executable(); strings.HasSuffix(xing_COM32_경로, 실행화일명) {
//			return 프로세스.Pid()
//		}
//	}
//
//	return -1
//}

func f접속유지_실행() {
	if !lib.F인터넷에_접속됨() {
		return
	}

	if 접속유지_실행중.G값() {
		return
	} else if 에러 := 접속유지_실행중.S값(true); 에러 != nil {
		return
	}

	go f접속유지_도우미()
}

func f접속유지_도우미() {
	defer 접속유지_실행중.S값(false)

	ch종료 := lib.F공통_종료_채널()
	에러_연속_발생_횟수 := 0

	for {
		lib.F대기(lib.P10초)

		if _, 에러 := (<-F시각_조회_t0167()).G값(); 에러 == nil {
			에러_연속_발생_횟수 = 0
		} else {
			에러_연속_발생_횟수++
		}

		// 3회 연속 에러 발생하면 API 연결에 문제 있다고 판단하고, C32 API 모듈 재시작.
		if 에러_연속_발생_횟수 >= 3 {
			C32_재시작()
			에러_연속_발생_횟수 = 0
		}

		select {
		case <-ch종료:
			return
		default:
		}
	}
}

func f에러_발생(TR코드, 코드, 내용 string) bool {
	switch TR코드 {
	case TR현물_정상_주문,
		TR시간_조회,
		TR현물_호가_조회,
		TR현물_시세_조회,
		TR현물_기간별_조회,
		TR현물_당일_전일_분틱_조회,
		TR_ETF_시간별_추이,
		TR기업정보_요약,
		TR재무순위_종합,
		TR현물_차트_틱,
		TR현물_차트_분,
		TR현물_차트_일주월,
		TR증시_주변_자금_추이,
		TR현물_종목_조회:
		return 코드 != "00000"
	case TR현물_정정_주문:
		return 코드 != "00131"
	case TR현물_취소_주문:
		return 코드 != "00156"
	default:
		// 에러 출력 지우지 말 것.
		panic(lib.New에러with출력("판별 불가능한 TR코드 : '%v'\n코드 : '%v'\n내용 : '%v'", TR코드, 코드, 내용))
	}
}


func f데이터_복원_이중_응답(대기_항목 *c32_콜백_대기_항목, 수신값 *lib.S바이트_변환) (에러 error) {
	완전값 := new(S이중_응답_일반형)

	if 대기_항목.대기값 != nil {
		대기값 := 대기_항목.대기값.(I이중_응답)

		if 대기값.G응답1() != nil {
			완전값.M응답1 = 대기값.G응답1()
		}

		if 대기값.G응답2() != nil {
			완전값.M응답2 = 대기값.G응답2()
		}
	}

	switch 변환값 := 수신값.S해석기(F바이트_변환값_해석).G해석값_단순형().(type) {
	case error:
		return 변환값
	case I이중_응답:
		if 변환값.G응답1() != nil {
			완전값.M응답1 = 변환값.G응답1()
		}

		if 변환값.G응답2() != nil {
			완전값.M응답2 = 변환값.G응답2()
		}
	case I이중_응답1:
		완전값.M응답1 = 변환값
	case I이중_응답2:
		완전값.M응답2 = 변환값
	default:
		panic(lib.New에러with출력("예상하지 못한 자료형 문자열 : '%v'", 수신값.G자료형_문자열()))
	}

	대기_항목.대기값 = 완전값

	if 완전값.M응답1 != nil && 완전값.M응답2 != nil {
		대기_항목.데이터_수신 = true
	} else {
		대기_항목.데이터_수신 = false
	}

	return nil
}

func f데이터_복원_반복_조회(대기_항목 *c32_콜백_대기_항목, 수신값 *lib.S바이트_변환) (에러 error) {
	defer lib.S예외처리{M에러: &에러}.S실행()

	완전값 := new(S헤더_반복값_일반형)

	if 대기_항목.대기값 != nil {
		대기값 := 대기_항목.대기값.(I헤더_반복값_TR데이터)

		if 대기값.G헤더_TR데이터() != nil {
			완전값.M헤더 = 대기값.G헤더_TR데이터()
		}

		if 대기값.G반복값_모음_TR데이터() != nil {
			완전값.M반복값_모음 = 대기값.G반복값_모음_TR데이터()
		}
	}

	switch 변환값 := 수신값.S해석기(F바이트_변환값_해석).G해석값_단순형().(type) {
	default:
		panic(lib.New에러with출력("예상하지 못한 자료형 : '%T' '%v'", 변환값, 수신값.G자료형_문자열()))
	case error:
		lib.F에러_출력(변환값.Error())
		return 변환값
	case I헤더_반복값_TR데이터:
		if 변환값.G헤더_TR데이터() != nil {
			완전값.M헤더 = 변환값.G헤더_TR데이터()
		}

		if 변환값.G반복값_모음_TR데이터() != nil {
			완전값.M반복값_모음 = 변환값.G반복값_모음_TR데이터()
		}
	case I헤더_TR데이터:
		완전값.M헤더 = 변환값
	case I반복값_모음_TR데이터:
		완전값.M반복값_모음 = 변환값
	}

	대기_항목.대기값 = 완전값

	if 완전값.M헤더 != nil && 완전값.M반복값_모음 != nil {
		대기_항목.데이터_수신 = true
	} else {
		대기_항목.데이터_수신 = false
	}

	//lib.F체크포인트(대기_항목.식별번호, 대기_항목.TR코드, 대기_항목.데이터_수신, 대기_항목.메시지_수신, 대기_항목.응답_완료, 대기_항목.에러)

	return nil
}

func f처리_가능한_TR코드(TR코드 string) bool {
	switch TR코드 {
	case TR현물_정상_주문, TR현물_정정_주문, TR현물_취소_주문,
		TR시간_조회, TR현물_호가_조회, TR현물_시세_조회,
		TR현물_기간별_조회, TR현물_당일_전일_분틱_조회, TR_ETF_시세_조회, TR_ETF_시간별_추이,
		TR기업정보_요약, TR재무순위_종합, TR현물_차트_틱, TR현물_차트_분, TR현물_차트_일주월,
		TR증시_주변_자금_추이, TR현물_종목_조회:
		return true
	default:
		lib.F문자열_출력("예상하지 못한 TR코드 : '%v'", TR코드)
		return false
	}
}

func C32_재시작() (에러 error) {
	defer lib.S예외처리{M에러: &에러}.S실행()

	lib.F문자열_출력("** C32 재시작 **")

	lib.F확인(C32_종료())
	lib.F확인(f초기화_xing_C32())
	lib.F조건부_패닉(!f초기화_작동_확인(), "초기화 작동 확인 실패.")
	lib.F확인(f종목모음_설정())
	lib.F확인(f전일_당일_설정())

	fmt.Println("** C32 재시작 완료     **")

	return nil
}

func f자료형_크기_비교_확인() (에러 error) {
	lib.S예외처리{M에러: &에러}.S실행()

	lib.F조건부_패닉(Sizeof_C_TR_DATA != C.sizeof_TR_DATA, "C.TR_DATA 크기 불일치 ")
	lib.F조건부_패닉(Sizeof_C_MSG_DATA != C.sizeof_MSG_DATA, "C.MSG_DATA 크기 불일치 ")
	lib.F조건부_패닉(Sizeof_C_REALTIME_DATA != C.sizeof_REALTIME_DATA, "C.REALTIME_DATA 크기 불일치 ")
	lib.F조건부_패닉(unsafe.Sizeof(TR_DATA{}) != unsafe.Sizeof(C.TR_DATA_UNPACKED{}), "TR_DATA 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(REALTIME_DATA{}) != unsafe.Sizeof(C.REALTIME_DATA_UNPACKED{}), "REALTIME_DATA_UNPACKED 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(MSG_DATA{}) != unsafe.Sizeof(C.MSG_DATA_UNPACKED{}), "MSG_DATA_UNPACKED 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00600InBlock1{}) != unsafe.Sizeof(C.CSPAT00600InBlock1{}), "CSPAT00600InBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00600OutBlock1{}) != unsafe.Sizeof(C.CSPAT00600OutBlock1{}), "CSPAT00600OutBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00600OutBlock2{}) != unsafe.Sizeof(C.CSPAT00600OutBlock2{}), "CSPAT00600OutBlock2 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00600OutBlock{}) != unsafe.Sizeof(C.CSPAT00600OutBlock{}), "CSPAT00600OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00700InBlock1{}) != unsafe.Sizeof(C.CSPAT00700InBlock1{}), "CSPAT00700InBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00700OutBlock1{}) != unsafe.Sizeof(C.CSPAT00700OutBlock1{}), "CSPAT00700OutBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00700OutBlock2{}) != unsafe.Sizeof(C.CSPAT00700OutBlock2{}), "CSPAT00700OutBlock2 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00700OutBlock{}) != unsafe.Sizeof(C.CSPAT00700OutBlock{}), "CSPAT00700OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00800InBlock1{}) != unsafe.Sizeof(C.CSPAT00800InBlock1{}), "CSPAT00800InBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00800OutBlock1{}) != unsafe.Sizeof(C.CSPAT00800OutBlock1{}), "CSPAT00800OutBlock1 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00800OutBlock2{}) != unsafe.Sizeof(C.CSPAT00800OutBlock2{}), "CSPAT00800OutBlock2 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(CSPAT00800OutBlock{}) != unsafe.Sizeof(C.CSPAT00800OutBlock{}), "CSPAT00800OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(SC0_OutBlock{}) != unsafe.Sizeof(C.SC0_OutBlock{}), "SC0_OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(SC1_OutBlock{}) != unsafe.Sizeof(C.SC1_OutBlock{}), "SC1_OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(SC2_OutBlock{}) != unsafe.Sizeof(C.SC2_OutBlock{}), "SC2_OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(SC3_OutBlock{}) != unsafe.Sizeof(C.SC3_OutBlock{}), "SC3_OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(SC4_OutBlock{}) != unsafe.Sizeof(C.SC4_OutBlock{}), "SC4_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T0167OutBlock{}) != unsafe.Sizeof(C.T0167OutBlock{}), "T0167OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1101InBlock{}) != unsafe.Sizeof(C.T1101InBlock{}), "T1101InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1101OutBlock{}) != unsafe.Sizeof(C.T1101OutBlock{}), "T1101OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1102InBlock{}) != unsafe.Sizeof(C.T1102InBlock{}), "T1102InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1102OutBlock{}) != unsafe.Sizeof(C.T1102OutBlock{}), "T1102OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1305InBlock{}) != unsafe.Sizeof(C.T1305InBlock{}), "T1305InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1305OutBlock{}) != unsafe.Sizeof(C.T1305OutBlock{}), "T1305OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1305OutBlock1{}) != unsafe.Sizeof(C.T1305OutBlock1{}), "T1305OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1310InBlock{}) != unsafe.Sizeof(C.T1310InBlock{}), "T1310InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1310OutBlock{}) != unsafe.Sizeof(C.T1310OutBlock{}), "T1310OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1310OutBlock1{}) != unsafe.Sizeof(C.T1310OutBlock1{}), "T1310OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1901InBlock{}) != unsafe.Sizeof(C.T1901InBlock{}), "T1901InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1901OutBlock{}) != unsafe.Sizeof(C.T1901OutBlock{}), "T1901OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T1902InBlock{}) != unsafe.Sizeof(C.T1902InBlock{}), "T1902InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1902OutBlock{}) != unsafe.Sizeof(C.T1902OutBlock{}), "T1902OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T1902OutBlock1{}) != unsafe.Sizeof(C.T1902OutBlock1{}), "T1902OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T3320InBlock{}) != unsafe.Sizeof(C.T3320InBlock{}), "T3320InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T3320OutBlock{}) != unsafe.Sizeof(C.T3320OutBlock{}), "T3320OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T3320OutBlock1{}) != unsafe.Sizeof(C.T3320OutBlock1{}), "T3320OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T8411InBlock{}) != unsafe.Sizeof(C.T8411InBlock{}), "T8411InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8411OutBlock{}) != unsafe.Sizeof(C.T8411OutBlock{}), "T8411OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8411OutBlock1{}) != unsafe.Sizeof(C.T8411OutBlock1{}), "T8411OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T8412InBlock{}) != unsafe.Sizeof(C.T8412InBlock{}), "T8412InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8412OutBlock{}) != unsafe.Sizeof(C.T8412OutBlock{}), "T8412OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8412OutBlock1{}) != unsafe.Sizeof(C.T8412OutBlock1{}), "T8412OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T8413InBlock{}) != unsafe.Sizeof(C.T8413InBlock{}), "T8413InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8413OutBlock{}) != unsafe.Sizeof(C.T8413OutBlock{}), "T8413OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8413OutBlock1{}) != unsafe.Sizeof(C.T8413OutBlock1{}), "T8413OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T8428InBlock{}) != unsafe.Sizeof(C.T8428InBlock{}), "T8428InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8428OutBlock{}) != unsafe.Sizeof(C.T8428OutBlock{}), "T8428OutBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8428OutBlock1{}) != unsafe.Sizeof(C.T8428OutBlock1{}), "T8428OutBlock1 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(T8436InBlock{}) != unsafe.Sizeof(C.T8436InBlock{}), "T8436InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(T8436OutBlock{}) != unsafe.Sizeof(C.T8436OutBlock{}), "T8436OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(H1_InBlock{}) != unsafe.Sizeof(C.H1_InBlock{}), "H1_InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(H1_OutBlock{}) != unsafe.Sizeof(C.H1_OutBlock{}), "H1_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(H2_InBlock{}) != unsafe.Sizeof(C.H2_InBlock{}), "H2_InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(H2_OutBlock{}) != unsafe.Sizeof(C.H2_OutBlock{}), "H2_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(S3_InBlock{}) != unsafe.Sizeof(C.S3_InBlock{}), "S3_InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(S3_OutBlock{}) != unsafe.Sizeof(C.S3_OutBlock{}), "S3_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(YS3InBlock{}) != unsafe.Sizeof(C.YS3InBlock{}), "YS3InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(YS3OutBlock{}) != unsafe.Sizeof(C.YS3OutBlock{}), "YS3OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(I5_InBlock{}) != unsafe.Sizeof(C.I5_InBlock{}), "I5_InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(I5_OutBlock{}) != unsafe.Sizeof(C.I5_OutBlock{}), "I5_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(VI_InBlock{}) != unsafe.Sizeof(C.VI_InBlock{}), "VI_InBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(VI_OutBlock{}) != unsafe.Sizeof(C.VI_OutBlock{}), "VI_OutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(DVIInBlock{}) != unsafe.Sizeof(C.DVIInBlock{}), "DVIInBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(DVIOutBlock{}) != unsafe.Sizeof(C.DVIOutBlock{}), "DVIOutBlock 크기 불일치")

	lib.F조건부_패닉(unsafe.Sizeof(JIFInBlock{}) != unsafe.Sizeof(C.JIFInBlock{}), "JIFInBlock 크기 불일치")
	lib.F조건부_패닉(unsafe.Sizeof(JIFOutBlock{}) != unsafe.Sizeof(C.JIFOutBlock{}), "JIFOutBlock 크기 불일치")

	return nil
}
