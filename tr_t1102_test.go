/* Copyright (C) 2015-2018 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

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

Copyright (C) 2015-2018년 UnHa Kim (unha.kim@kuh.pe.kr)

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

import (
	"github.com/ghts/lib"
	"github.com/ghts/xing_types"
	"strings"
	"testing"
	"time"
)

func TestF현물_시세_조회_t1102(t *testing.T) {
	접속됨, 에러 := F접속됨()
	lib.F테스트_에러없음(t, 에러)
	lib.F테스트_참임(t, 접속됨)

	종목코드_모음 := []string{
		lib.F임의_종목_코스피_주식().G코드(),
		lib.F임의_종목_코스닥_주식().G코드(),
		lib.F임의_종목_ETF().G코드()}

	for _, 종목코드 := range 종목코드_모음 {
		f현물_시세조회_테스트_도우미(t, 종목코드)

		if t.Failed() {
			break
		}
	}
}

func f현물_시세조회_테스트_도우미(t *testing.T, 종목코드 string) {
	질의값 := lib.New질의값_단일종목()
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = xt.TR현물_시세_조회
	질의값.M종목코드 = 종목코드

	응답 := lib.New소켓_질의_단순형(lib.P주소_Xing_TR, lib.F임의_변환_형식(), lib.P1분).S질의(질의값).G응답_검사()

	ETF종목_여부 := etf종목_여부(질의값.M종목코드)

	값 := new(xt.S현물시세조회_응답)
	lib.F테스트_에러없음(t, 응답.G값(0, 값))
	lib.F테스트_다름(t, 값.M한글명, "")
	lib.F테스트_참임(t, 값.M현재가 >= 0)
	lib.F테스트_같음(t, 값.M전일대비구분, xt.P구분_상한, xt.P구분_상승, xt.P구분_보합, xt.P구분_하한, xt.P구분_하락)
	lib.F테스트_참임(t, 값.M전일대비등락폭 >= 0)

	switch 값.M전일대비구분 { // 등락율 확인
	case xt.P구분_상한, xt.P구분_상승:
		lib.F테스트_참임(t, 값.M등락율 >= 0)
	case xt.P구분_하한, xt.P구분_하락:
		lib.F테스트_참임(t, 값.M등락율 <= 0)
	case xt.P구분_보합:
		lib.F테스트_같음(t, 값.M등락율, 0)
	}

	lib.F테스트_참임(t, 값.M누적거래량 >= 0)
	lib.F메모("현물 시세조회의 가중평균이 무슨 값이지?")
	//lib.F문자열_출력("가중평균 : %v", 값.M가중평균)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M현재가)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M하한가)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M시가)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M고가)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M저가)
	lib.F테스트_참임(t, 값.M52주_최고가 >= 값.M기준가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M현재가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M하한가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M시가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M고가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M저가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M기준가)
	lib.F테스트_참임(t, 값.M상한가 >= 값.M52주_최저가)
	lib.F테스트_참임(t, 값.M하한가 <= 값.M현재가)
	lib.F테스트_참임(t, 값.M하한가 <= 값.M시가, 종목코드, 값.M하한가, 값.M시가)
	lib.F테스트_참임(t, 값.M하한가 <= 값.M고가)
	lib.F테스트_참임(t, 값.M하한가 <= 값.M저가)
	lib.F테스트_참임(t, 값.M하한가 <= 값.M기준가)
	lib.F테스트_참임(t, 값.M고가 >= 값.M현재가)
	lib.F테스트_참임(t, 값.M고가 >= 값.M시가)
	lib.F테스트_참임(t, 값.M고가 >= 값.M저가)
	lib.F테스트_참임(t, 값.M고가 >= 값.M52주_최저가)
	lib.F테스트_참임(t, 값.M저가 <= 값.M현재가)
	lib.F테스트_참임(t, 값.M저가 >= 값.M52주_최저가)
	lib.F테스트_참임(t, 값.M52주_최저가 > 0)
	lib.F테스트_참임(t, 값.M전일거래량 >= 0)
	lib.F테스트_같음(t, 값.M거래량차, lib.F절대값_실수(값.M전일거래량-값.M누적거래량))

	현재 := time.Now()
	금일_자정 := time.Date(현재.Year(), 현재.Month(), 현재.Day(), 0, 0, 0, 0, 현재.Location())
	개장_시간 := time.Date(현재.Year(), 현재.Month(), 현재.Day(), 8, 0, 0, 0, 현재.Location())
	lib.F테스트_참임(t, 값.M시가시간.After(개장_시간.Add(-1*lib.P10분)))
	lib.F테스트_참임(t, 값.M시가시간.Before(현재.Add(lib.P10분)))
	lib.F테스트_참임(t, 값.M고가시간.After(개장_시간.Add(-1*lib.P10분)))
	lib.F테스트_참임(t, 값.M고가시간.Before(현재.Add(lib.P10분)))
	lib.F테스트_참임(t, 값.M고가시간.After(값.M시가시간) || 값.M고가시간.Equal(값.M시가시간))
	lib.F테스트_참임(t, 값.M저가시간.After(개장_시간.Add(-1*lib.P10분)))
	lib.F테스트_참임(t, 값.M저가시간.Before(현재.Add(lib.P10분)))
	lib.F테스트_참임(t, 값.M저가시간.After(값.M시가시간) || 값.M저가시간.Equal(값.M시가시간))
	lib.F테스트_참임(t, 값.M52주_최고가일.After(금일_자정.Add(-380*lib.P1일))) // 영업일 기준이므로, 휴일을 고려해서 여유를 둠.
	lib.F테스트_참임(t, 값.M52주_최고가일.Before(금일_자정.Add(lib.P1일)))
	lib.F테스트_참임(t, 값.M52주_최저가일.After(금일_자정.Add(-380*lib.P1일)),
		값.M52주_최저가일, 금일_자정.Add(-62*lib.P1일)) // 영업일 기준이므로, 휴일을 고려해서 여유를 둠.
	lib.F테스트_참임(t, 값.M52주_최저가일.Before(금일_자정.Add(lib.P1일)))
	lib.F테스트_참임(t, 값.M소진율 >= 0.0 && 값.M소진율 <= 100.0, 값.M소진율)
	//lib.F문자열_출력("값.PER : %v", 값.PER)
	//lib.F문자열_출력("값.PBRX : %v", 값.PBRX)
	lib.F테스트_참임(t, 값.M상장주식수_천 > 0, 값.M상장주식수_천)
	lib.F테스트_참임(t, 값.M증거금율 >= 0 && 값.M증거금율 <= 100, 값.M증거금율)
	lib.F테스트_같음(t, 값.M수량단위, 1) // 단주 거래 전면 허용되었다고 하던 데.
	lib.F테스트_같음(t, len(값.M매도증권사코드_모음), 5)
	lib.F테스트_같음(t, len(값.M매수증권사코드_모음), 5)
	lib.F테스트_같음(t, len(값.M매도증권사명_모음), 5)
	lib.F테스트_같음(t, len(값.M매수증권사명_모음), 5)
	lib.F테스트_같음(t, len(값.M총매도수량_모음), 5)
	lib.F테스트_같음(t, len(값.M총매수수량_모음), 5)
	lib.F테스트_같음(t, len(값.M매도증감_모음), 5)
	lib.F테스트_같음(t, len(값.M매수증감_모음), 5)
	lib.F테스트_같음(t, len(값.M매도비율_모음), 5)
	lib.F테스트_같음(t, len(값.M매수비율_모음), 5)

	for i := range 값.M매도증권사코드_모음 {
		if 값.M매도증권사코드_모음[i] == "" {
			lib.F테스트_같음(t, strings.TrimSpace(값.M매도증권사명_모음[i]), "")
			lib.F테스트_같음(t, 값.M총매도수량_모음[i], 0)
			lib.F테스트_같음(t, 값.M총매수수량_모음[i], 0)
			lib.F테스트_같음(t, 값.M매도증감_모음[i], 0)
			lib.F테스트_같음(t, 값.M매수증감_모음[i], 0)
			lib.F테스트_같음(t, 값.M매도비율_모음[i], 0)
			lib.F테스트_같음(t, 값.M매수비율_모음[i], 0)
		} else {
			lib.F테스트_다름(t, 값.M매도증권사명_모음[i], "")
			lib.F테스트_다름(t, 값.M총매도수량_모음[i], 0)
			lib.F테스트_다름(t, 값.M총매수수량_모음[i], 0)
			lib.F테스트_다름(t, 값.M매도증감_모음[i], 0)
			lib.F테스트_다름(t, 값.M매수증감_모음[i], 0)
			lib.F테스트_다름(t, 값.M매도비율_모음[i], 0)
			lib.F테스트_다름(t, 값.M매수비율_모음[i], 0)
		}
	}

	lib.F테스트_참임(t, 값.M외국계_매도_합계수량 >= 0)
	lib.F테스트_참임(t, 값.M외국계_매도_직전대비 >= 0, 값.M외국계_매도_직전대비)
	lib.F테스트_참임(t, 값.M외국계_매도_비율 >= 0 && 값.M외국계_매도_비율 <= 100, 값.M외국계_매도_비율)
	lib.F테스트_참임(t, 값.M외국계_매수_합계수량 >= 0)
	lib.F테스트_참임(t, 값.M외국계_매수_직전대비 >= 0, 값.M외국계_매수_직전대비)
	lib.F테스트_참임(t, 값.M외국계_매수_비율 >= 0 && 값.M외국계_매수_비율 <= 100, 값.M외국계_매수_비율)
	lib.F테스트_참임(t, 값.M회전율 >= 0, 값.M회전율)
	lib.F테스트_같음(t, len(값.M종목코드), 6)
	lib.F테스트_참임(t, 값.M누적거래대금 >= 0, 값.M누적거래대금)
	lib.F테스트_참임(t, 값.M전일동시간거래량 >= 0, 값.M전일동시간거래량)
	lib.F테스트_참임(t, 값.M연중_최고가 > 0)
	lib.F테스트_같음(t, 값.M연중_최고가_일자.Year(), time.Now().Year())
	lib.F테스트_참임(t, 값.M연중_최저가 > 0)
	lib.F테스트_같음(t, 값.M연중_최저가_일자.Year(), time.Now().Year())
	lib.F테스트_참임(t, 값.M목표가 >= 0, 값.M목표가)

	if ETF종목_여부 {
		lib.F테스트_같음(t, 값.M자본금, 0)
		lib.F테스트_같음(t, 값.M액면가, 0)
		lib.F테스트_같음(t, 값.M전분기_매출액, 0)
		lib.F테스트_같음(t, 값.M전전분기_매출액, 0)
		lib.F테스트_같음(t, 값.M전년대비매출액, 0)
		//lib.F테스트_같음(t, 값.M발행가격, 0)
		lib.F테스트_같음(t, 값.M결산월, 0)
	} else {
		lib.F테스트_참임(t, 값.M자본금 >= 0, 종목코드, ETF종목_여부, 값.M자본금)
		lib.F테스트_참임(t, 값.M액면가 >= 0, 종목코드, ETF종목_여부, 값.M액면가) // 완리 (900180)는 액면가 0임
		lib.F테스트_참임(t, 값.M전분기_매출액 >= 0)                     // '동북아 12호(083370)'의 경우 대여업만 하므로 (판매) 매출액 0임.
		lib.F테스트_참임(t, 값.M전전분기_매출액 >= 0)
		//lib.F테스트_참임(t, 값.M전년대비매출액 > 0)
		//lib.F테스트_참임(t, 값.M발행가격 >= 0, 종목코드, 값.M발행가격)   // 발행가격이 0인 종목이 존재함. 이해불가.
		lib.F테스트_같음(t, 값.M결산월, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	}

	lib.F테스트_참임(t, 값.M유동주식수 > 0)

	오차 := lib.F절대값_실수((float64(값.M대용가) - float64(값.M기준가)*float64(값.M증거금율)/100) / float64(값.M대용가))
	lib.F테스트_참임(t, 오차 < 3, "오차:%v, 기준가:%v, 증거금율:%v, 대용가:%v", 오차, 값.M기준가, 값.M증거금율, 값.M대용가)

	시가총액 := 값.M시가총액_억 * 100000000
	상장주식수 := 값.M상장주식수_천 * 1000
	오차율 := lib.F절대값_실수((시가총액-값.M현재가*상장주식수)/시가총액) / 100
	lib.F테스트_참임(t, 오차율 < 10, 오차율)
	lib.F테스트_참임(t, 값.M상장일.After(time.Time{})) // 1900-01-01 이후
	lib.F테스트_참임(t, 값.M전분기명 == "" ||
		lib.F정규식_검색(값.M전분기명, []string{"[0-9]+ [1-4]분기"}) != "" ||
		lib.F정규식_검색(값.M전분기명, []string{"[0-9]+ 결산"}) != "", 값.M전분기명)

	//lib.F테스트_참임(t, 값.M전분기_영업이익 ???)
	//lib.F테스트_참임(t, 값.M전분기_경상이익 ???)
	//lib.F테스트_참임(t, 값.M전분기_순이익 ???)
	lib.F테스트_참임(t, float64(값.M전분기_순이익)*값.M전분기EPS >= 0, 값.M전분기_순이익, 값.M전분기EPS)
	lib.F테스트_참임(t, 값.M전전분기명 == "" || lib.F정규식_검색(값.M전전분기명, []string{"[0-9]+ [1-4]분기"}) != "", 값.M전전분기명)

	//lib.F테스트_참임(t, 값.M전전분기_영업이익 ???)
	//lib.F테스트_참임(t, 값.M전전분기_경상이익 ???)
	//lib.F테스트_참임(t, 값.M전전분기_순이익 ???)
	lib.F테스트_참임(t, float64(값.M전전분기_순이익)*값.M전전분기EPS >= 0, 값.M전전분기_순이익, 값.M전전분기EPS)

	//lib.F테스트_참임(t, 값.M전년대비영업이익 ???)
	//lib.F테스트_참임(t, 값.M전년대비경상이익 ???)
	//lib.F테스트_참임(t, 값.M전년대비순이익 ???)

	// 주식소각, 주식분할등 주식수량에 변동이 있는 경우, EPS와 순이익의 방향성이 다를 수 있다.
	//lib.F테스트_참임(t, float64(값.M전년대비순이익) * 값.M전년대비EPS >= 0, 종목코드, 값.M전년대비순이익, 값.M전년대비EPS)

	//lib.F변수값_확인(값.M락구분)      // ??
	//lib.F변수값_확인(값.M관리_급등구분)  // ??
	//lib.F변수값_확인(값.M정지_연장구분)  // ??
	//lib.F변수값_확인(값.M투자_불성실구분) // ??
	lib.F테스트_같음(t, 값.M시장구분, lib.P시장구분_전체, lib.P시장구분_코스피, lib.P시장구분_코스닥)
	//lib.F변수값_확인(값.T_PER) // ??
	lib.F테스트_같음(t, 값.M통화ISO코드, "KRW")
	lib.F테스트_같음(t, len(값.M총매도대금_모음), 5)
	lib.F테스트_같음(t, len(값.M총매수대금_모음), 5)
	lib.F테스트_같음(t, len(값.M총매도평단가_모음), 5)
	lib.F테스트_같음(t, len(값.M총매수평단가_모음), 5)

	for i := range 값.M총매도대금_모음 {
		lib.F테스트_참임(t, 값.M총매도대금_모음[i] >= 0)
		lib.F테스트_참임(t, 값.M총매수대금_모음[i] >= 0)
		lib.F테스트_참임(t, 값.M총매도평단가_모음[i] >= 0)
		lib.F테스트_참임(t, 값.M총매수평단가_모음[i] >= 0)
	}

	lib.F테스트_참임(t, 값.M외국계매도대금 >= 0)
	lib.F테스트_참임(t, 값.M외국계매수대금 >= 0)
	lib.F테스트_참임(t, 값.M외국계매도평단가 >= 0)
	lib.F테스트_참임(t, 값.M외국계매수평단가 >= 0)

	//lib.F변수값_확인(값.M투자주의환기)     // ??
	//lib.F변수값_확인(값.M기업인수목적회사여부) // ??
	//lib.F변수값_확인(값.M배분적용구분코드)   // ??
	//lib.F변수값_확인(값.M배분적용구분)     // ??
	//lib.F변수값_확인(값.M단기과열_VI발동)  // ??
	//lib.F변수값_확인(값.M정적VI상한가)    // ??
	//lib.F변수값_확인(값.M정적VI하한가)    // ??
}
