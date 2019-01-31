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
(자유 소프트웨어 재단 : Free Software Foundation, In,
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

import (
	"github.com/ghts/lib"
	"time"

	"sync"
)

var (
	소켓REP_TR콜백 lib.I소켓Raw
	소켓REQ_저장소  = lib.New소켓_저장소(20, func() lib.I소켓_질의 {
		return lib.NewNano소켓REQ_단순형(lib.P주소_Xing_C함수_호출, lib.P30초)
	})
	소켓SUB_실시간_정보 lib.I소켓Raw

	ch질의 = make(chan *lib.S작업, 1000)
	ch신호_C32_모음 []chan T신호_C32

	대기소_C32 = newC32_콜백_대기_저장소()

	tr코드별_10분당_전송_제한 = make(map[string]lib.I전송_권한)
	tr코드별_초당_전송_제한   = make(map[string]lib.I전송_권한)

	전일, 당일    lib.I안전한_시각
	최근_영업일_모음 []time.Time

	xing_C32_실행_잠금 sync.Mutex
	xing_C32_경로    = lib.F_GOPATH() + `/src/github.com/ghts/xing_C32/xing_C32.bat`

	//xing_COM32_실행_잠금 sync.Mutex
	//xing_COM32_경로    = lib.F_GOPATH() + `/src/github.com/ghts/xing_COM32/xing_COM32.bat`

	접속유지_실행중   = lib.New안전한_bool(false)
	주문_응답_구독_중 = lib.New안전한_bool(false)
)

// 종목 관련 저장소는 초기화 이후에는 사실상 읽기 전용. 다중 사용에 문제가 없음.
var (
	종목모음_설정_잠금 sync.Mutex
	종목모음_설정일   = lib.New안전한_시각(time.Time{})
	종목맵_전체     = make(map[string]*lib.S종목)
	종목모음_코스피   = make([]*lib.S종목, 0)
	종목모음_코스닥   = make([]*lib.S종목, 0)
	종목모음_ETF   = make([]*lib.S종목, 0)
	종목모음_ETN   = make([]*lib.S종목, 0)
	종목모음_전체    = make([]*lib.S종목, 0)
	기준가_맵      = make(map[string]int64)
	하한가_맵      = make(map[string]int64)
)
