package tencentcos

import "testing"

func Test_sig(t *testing.T) {
	s := sig{
		appID:       "200001",
		bucket:      "newbucket",
		secretID:    "AKIDUfLUEUigQiXqm7CVSspKJnuaiIKtxqAv",
		expiredTime: "1470737000",
		currentTime: "1470736940",
		rand:        "490258943",
		fileid:      "",
	}
	ss := s.getMultiEffectSignature("bLcPnl88WU30VY57ipRhSePfPdOfSruK")
	expect := "v6+um3VE3lxGz97PmnSg6+/V9PZhPTIwMDAwMSZiPW5ld2J1Y2tldCZrPUFLSURVZkxVRVVpZ1FpWHFtN0NWU3NwS0pudWFpSUt0eHFBdiZlPTE0NzA3MzcwMDAmdD0xNDcwNzM2OTQwJnI9NDkwMjU4OTQzJmY9"
	if ss != expect {
		t.Error("expect:", expect, "\n", "result:", ss)
	}

	s = sig{
		appID:       "200001",
		bucket:      "newbucket",
		secretID:    "AKIDUfLUEUigQiXqm7CVSspKJnuaiIKtxqAv",
		expiredTime: "0",
		currentTime: "1470736940",
		rand:        "490258943",
		fileid:      "/200001/newbucket/tencent_test.jpg",
	}
	ss = s.getOnceSignature("bLcPnl88WU30VY57ipRhSePfPdOfSruK")
	expect = "CkZ0/gWkHy3f76ER7k6yXgzq7w1hPTIwMDAwMSZiPW5ld2J1Y2tldCZrPUFLSURVZkxVRVVpZ1FpWHFtN0NWU3NwS0pudWFpSUt0eHFBdiZlPTAmdD0xNDcwNzM2OTQwJnI9NDkwMjU4OTQzJmY9LzIwMDAwMS9uZXdidWNrZXQvdGVuY2VudF90ZXN0LmpwZw=="
	if ss != expect {
		t.Error("expect:", expect, "\n", "result:", ss)
	}
}
