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

import "C"
import (
	"bytes"
	"encoding/binary"
	"github.com/ghts/lib"
	"strings"
	"time"
)

func F현물_체결_미체결_조회_t0425(계좌번호, 비밀번호, 종목코드 string, 체결_구분 lib.T체결_구분,
	매도_매수_구분 lib.T매도_매수_구분) (응답값_모음 []*S현물_체결_미체결_조회_응답_반복값_t0425, 에러 error) {
	defer lib.S예외처리{M에러: &에러, M함수: func() { 응답값_모음 = nil }}.S실행()

	lib.F확인(F종목코드_검사(종목코드))

	응답값_모음 = make([]*S현물_체결_미체결_조회_응답_반복값_t0425, 0)
	연속키 := ""

	F접속_확인()

	for {
		질의값 := new(S질의값_체결_미체결_조회_t0425)
		질의값.S질의값_기본형 = lib.New질의값_기본형(TR조회, TR체결_미체결_조회_t0425)
		질의값.M계좌번호 = 계좌번호
		질의값.M비밀번호 = 비밀번호
		질의값.M종목코드 = 종목코드
		질의값.M체결구분 = 체결_구분
		질의값.M매도_매수_구분 = 매도_매수_구분
		질의값.M정렬순서 = P주문번호_순
		//질의값.M정렬순서 = P주문번호_역순
		질의값.M연속키 = 연속키

		i응답값, 에러 := F질의_단일TR(질의값)
		lib.F확인(에러)

		값, ok := i응답값.(*S현물_체결_미체결_조회_응답_t0425)
		lib.F조건부_패닉(!ok, "예상하지 못한 자료형 : '%T'", i응답값)

		연속키 = 값.M헤더.M연속키
		응답값_모음 = append(값.M반복값_모음, 응답값_모음...)

		if lib.F2문자열_공백제거(연속키) == "" {
			break
		}
	}

	return 응답값_모음, nil
}

type S질의값_체결_미체결_조회_t0425 struct {
	*lib.S질의값_기본형
	M계좌번호     string
	M비밀번호     string
	M종목코드     string
	M체결구분     lib.T체결_구분
	M매도_매수_구분 lib.T매도_매수_구분
	M정렬순서     T정렬_순서_t0425
	M연속키      string
}

type S현물_체결_미체결_조회_응답_t0425 struct {
	M헤더     *S현물_체결_미체결_조회_응답_헤더_t0425
	M반복값_모음 []*S현물_체결_미체결_조회_응답_반복값_t0425
}

type S현물_체결_미체결_조회_응답_헤더_t0425 struct {
	M총_주문_수량    int64
	M총_체결_수량    int64
	M총_미체결_수량   int64
	M추정_수수료     int64
	M총_주문_금액    int64
	M총_매도_체결_금액 int64
	M총_매수_체결_금액 int64
	M추정_제세금     int64
	M연속키        string
}

type S현물_체결_미체결_조회_응답_반복값_t0425 struct {
	M주문_번호   int64
	M종목코드    string
	M매매_구분   string
	M주문_수량   int64
	M주문_가격   int64
	M체결_수량   int64
	M체결_가격   int64
	M미체결_잔량  int64
	M확인_수량   int64
	M상태      string
	M원_주문_번호 int64
	M유형      string
	M주문_시간   time.Time
	M주문_매체   string
	M처리_순번   int64
	M호가_유형   T호가_유형_t0425
	M현재가     int64
	M주문_구분   T주문_구분_t0425
	M신용_구분   T신용_구분_t0425
	M대출_일자   time.Time
}

func NewT0425InBlock(질의값 *S질의값_체결_미체결_조회_t0425) (g *T0425InBlock) {
	g = new(T0425InBlock)

	lib.F바이트_복사_문자열(g.Accno[:], 질의값.M계좌번호)
	lib.F바이트_복사_문자열(g.Passwd[:], 질의값.M비밀번호)
	lib.F바이트_복사_문자열(g.Expcode[:], 질의값.M종목코드)
	lib.F바이트_복사_문자열(g.Chegb[:], lib.F2문자열(int(질의값.M체결구분)))
	lib.F바이트_복사_문자열(g.Medosu[:], lib.F2문자열(int(질의값.M매도_매수_구분)))
	lib.F바이트_복사_문자열(g.Sortgb[:], lib.F2문자열(int(질의값.M정렬순서)))
	lib.F바이트_복사_문자열(g.Ordno[:], 질의값.M연속키)

	f속성값_초기화(g)

	return g
}

func New현물_체결_미체결_조회_응답_t0425(b []byte) (s *S현물_체결_미체결_조회_응답_t0425, 에러 error) {
	defer lib.S예외처리{M에러: &에러, M함수: func() { s = nil }}.S실행()

	lib.F조건부_패닉(len(b) < SizeT0425OutBlock+5, "예상하지 못한 길이 : '%v", len(b))
	lib.F조건부_패닉((len(b)-(SizeT0425OutBlock+5))%SizeT0425OutBlock1 != 0, "예상하지 못한 길이 : '%v", len(b))

	s = new(S현물_체결_미체결_조회_응답_t0425)
	s.M헤더, 에러 = new현물_체결_미체결_조회_응답_헤더_t0425(b[:SizeT0425OutBlock])
	lib.F확인(에러)

	s.M반복값_모음, 에러 = new현물_체결_미체결_조회_응답_반복값_모음_t0425(b[SizeT0425OutBlock+5:])
	lib.F확인(에러)

	return s, nil
}

func new현물_체결_미체결_조회_응답_헤더_t0425(b []byte) (s *S현물_체결_미체결_조회_응답_헤더_t0425, 에러 error) {
	defer lib.S예외처리{M에러: &에러, M함수: func() { s = nil }}.S실행()

	lib.F조건부_패닉(len(b) != SizeT0425OutBlock, "예상하지 못한 길이 : '%v", len(b))

	g := new(T0425OutBlock)
	lib.F확인(binary.Read(bytes.NewBuffer(b), binary.BigEndian, g))

	s = new(S현물_체결_미체결_조회_응답_헤더_t0425)
	s.M총_주문_수량 = lib.F2정수64_단순형(g.Tqty)
	s.M총_체결_수량 = lib.F2정수64_단순형(g.Tcheqty)
	s.M총_미체결_수량 = lib.F2정수64_단순형(g.Tordrem)
	s.M추정_수수료 = lib.F2정수64_단순형(g.Cmss)
	s.M총_주문_금액 = lib.F2정수64_단순형(g.Tamt)
	s.M총_매도_체결_금액 = lib.F2정수64_단순형(g.Tmdamt)
	s.M총_매수_체결_금액 = lib.F2정수64_단순형(g.Tmsamt)
	s.M추정_제세금 = lib.F2정수64_단순형(g.Tax)
	s.M연속키 = lib.F2문자열(g.Ordno)

	return s, nil
}

func new현물_체결_미체결_조회_응답_반복값_모음_t0425(b []byte) (값_모음 []*S현물_체결_미체결_조회_응답_반복값_t0425, 에러 error) {
	defer lib.S예외처리{M에러: &에러}.S실행()

	나머지 := len(b) % SizeT0425OutBlock1
	lib.F조건부_패닉(나머지 != 0, "예상하지 못한 길이. '%v' '%v'", len(b), 나머지)

	버퍼 := bytes.NewBuffer(b)
	수량 := len(b) / SizeT0425OutBlock1
	g_모음 := make([]*T0425OutBlock1, 수량, 수량)
	값_모음 = make([]*S현물_체결_미체결_조회_응답_반복값_t0425, 수량, 수량)

	for i, g := range g_모음 {
		g = new(T0425OutBlock1)
		lib.F확인(binary.Read(버퍼, binary.BigEndian, g))

		주문시간_문자열 := string(g.Ordtime[:])
		주문시간_문자열 = 주문시간_문자열[:6] + "." + 주문시간_문자열[6:]

		값 := new(S현물_체결_미체결_조회_응답_반복값_t0425)
		값.M주문_번호 = lib.F2정수64_단순형(g.Ordno)
		값.M종목코드 = lib.F2문자열_공백제거(g.Expcode)
		값.M매매_구분 = lib.F2문자열_EUC_KR_공백제거(g.Medosu)
		값.M주문_수량 = lib.F2정수64_단순형(g.Qty)
		값.M주문_가격 = lib.F2정수64_단순형(g.Price)
		값.M체결_수량 = lib.F2정수64_단순형(g.Cheqty)
		값.M체결_가격 = lib.F2정수64_단순형(g.Cheprice)
		값.M미체결_잔량 = lib.F2정수64_단순형(g.Ordrem)
		값.M확인_수량 = lib.F2정수64_단순형(g.Cfmqty)
		값.M상태 = lib.F2문자열_EUC_KR_공백제거(g.Status)
		값.M원_주문_번호 = lib.F2정수64_단순형_공백은_0(g.Orgordno)
		값.M유형 = lib.F2문자열_EUC_KR_공백제거(g.Ordgb)
		값.M주문_시간 = F2당일_시각_단순형("150405.99", 주문시간_문자열)
		값.M주문_매체 = lib.F2문자열_EUC_KR_공백제거(g.Ordermtd)
		값.M처리_순번 = lib.F2정수64_단순형(g.Sysprocseq)
		값.M호가_유형 = T호가_유형_t0425(lib.F2정수64_단순형_공백은_0(g.Hogagb))
		값.M현재가 = lib.F2정수64_단순형(g.Price1)
		값.M주문_구분 = T주문_구분_t0425(lib.F2문자열_공백제거(g.Orggb))
		값.M신용_구분 = T신용_구분_t0425(lib.F2정수64_단순형_공백은_0(g.Singb))
		값.M대출_일자 = lib.F2포맷된_일자_단순형_공백은_초기값("20060102", g.Loandt)

		if len(값.M종목코드) == 7 && strings.HasPrefix(값.M종목코드, "A") {
			값.M종목코드 = 값.M종목코드[1:]
		}

		값_모음[i] = 값
	}

	return 값_모음, nil
}
