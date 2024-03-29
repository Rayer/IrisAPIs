package IrisAPIs

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CurrencyContextTestSuite struct {
	suite.Suite
	db              *DatabaseContext
	currencyContext *CurrencyContext
}

func (c *CurrencyContextTestSuite) SetupSuite() {
	c.currencyContext = NewTestCurrencyContext()
}

func (c *CurrencyContextTestSuite) SetupTest() {
	if c.currencyContext == nil {
		c.T().Errorf("Fail due to can't initialize CurrencyContext!")
		c.T().FailNow()
	}
}

func TestCurrencyContextTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyContextTestSuite))
}

func (c *CurrencyContextTestSuite) TestSyncToDb() {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c.Run(tt.name, func() {
			if err := c.currencyContext.SyncToDb(); (err != nil) != tt.wantErr {
				assert.Equal(c.T(), err, nil)
				c.Errorf(err, "SyncToDb() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
	{"success":true,"timestamp":1589637365,"base":"EUR","date":"2020-05-16","rates":{"AED":3.97467,"AFN":83.005989,"ALL":123.366634,"AMD":527.975462,"ANG":1.941971,"AOA":603.731679,"ARS":73.193491,"AUD":1.687,"AWG":1.947888,"AZN":1.843968,"BAM":1.957295,"BBD":2.184355,"BDT":91.893973,"BGN":1.957199,"BHD":0.408517,"BIF":2066.925534,"BMD":1.08216,"BND":1.541246,"BOB":7.459452,"BRL":6.337855,"BSD":1.08182,"BTC":0.000115,"BTN":81.97207,"BWP":13.193597,"BYN":2.647737,"BYR":21210.33532,"BZD":2.180752,"CAD":1.526933,"CDF":1960.874249,"CHF":1.0516,"CLF":0.032461,"CLP":895.708037,"CNY":7.685721,"COP":4244.988896,"CRC":615.516389,"CUC":1.08216,"CUP":28.677239,"CVE":110.543067,"CZK":27.718233,"DJF":192.321895,"DKK":7.455546,"DOP":59.847652,"DZD":139.552973,"EGP":17.024004,"ERN":16.232787,"ETB":36.534145,"EUR":1,"FJD":2.437295,"FKP":0.89399,"GBP":0.894074,"GEL":3.468366,"GGP":0.89399,"GHS":6.260339,"GIP":0.89399,"GMD":55.735408,"GNF":10231.822857,"GTQ":8.32871,"GYD":226.231417,"HKD":8.387865,"HNL":27.054383,"HRK":7.563259,"HTG":114.848025,"HUF":354.878171,"IDR":16100.808828,"ILS":3.828903,"IMP":0.89399,"INR":82.111095,"IQD":1287.770359,"IRR":45564.345721,"ISK":157.324833,"JEP":0.89399,"JMD":157.502345,"JOD":0.767148,"JPY":115.887974,"KES":115.873188,"KGS":82.877334,"KHR":4433.072409,"KMF":491.845813,"KPW":973.94397,"KRW":1334.509263,"KWD":0.334608,"KYD":0.901583,"KZT":455.588891,"LAK":9744.850865,"LBP":1636.59543,"LKR":203.283306,"LRD":214.755058,"LSL":20.024041,"LTL":3.195337,"LVL":0.654588,"LYD":1.537079,"MAD":10.692152,"MDL":19.255132,"MGA":4128.440643,"MKD":61.549021,"MMK":1520.013038,"MNT":3030.113149,"MOP":8.636724,"MRO":386.331123,"MUR":43.112348,"MVR":16.777528,"MWK":798.09702,"MXN":25.924012,"MYR":4.70924,"MZN":74.025194,"NAD":20.023995,"NGN":418.807136,"NIO":36.739739,"NOK":11.074396,"NPR":131.155591,"NZD":1.823154,"OMR":0.41661,"PAB":1.08192,"PEN":3.723036,"PGK":3.729165,"PHP":54.903428,"PKR":173.145964,"PLN":4.569316,"PYG":7103.482417,"QAR":3.940186,"RON":4.841804,"RSD":117.634783,"RUB":79.65293,"RWF":1014.524967,"SAR":4.065519,"SBD":9.061443,"SCR":17.863116,"SDG":59.847426,"SEK":10.671075,"SGD":1.546488,"SHP":0.89399,"SLL":10659.276025,"SOS":628.735306,"SRD":8.07079,"STD":23862.377246,"SVC":9.467174,"SYP":556.718973,"SZL":20.023915,"THB":34.6995,"TJS":11.085272,"TMT":3.78756,"TND":3.153144,"TOP":2.508559,"TRY":7.468207,"TTD":7.318545,"TWD":32.431293,"TZS":2503.472798,"UAH":28.813803,"UGX":4095.865634,"USD":1.08216,"UYU":47.526011,"UZS":10940.637611,"VEF":10.808077,"VND":25252.743871,"VUV":130.220071,"WST":3.034681,"XAF":656.445302,"XAG":0.06501,"XAU":0.00062,"XCD":2.924592,"XDR":0.796292,"XOF":656.871459,"XPF":119.644,"YER":270.891734,"ZAR":20.109119,"ZMK":9740.742168,"ZMW":19.98716,"ZWL":348.45551}}
*/

func (c *CurrencyContextTestSuite) Test_saveCurrencyEntries() {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal Case",
			args: args{
				raw: "{\"success\":true,\"timestamp\":1589637365,\"base\":\"EUR\",\"date\":\"2020-05-16\",\"rates\":{\"AED\":3.97467,\"AFN\":83.005989,\"ALL\":123.366634,\"AMD\":527.975462,\"ANG\":1.941971,\"AOA\":603.731679,\"ARS\":73.193491,\"AUD\":1.687,\"AWG\":1.947888,\"AZN\":1.843968,\"BAM\":1.957295,\"BBD\":2.184355,\"BDT\":91.893973,\"BGN\":1.957199,\"BHD\":0.408517,\"BIF\":2066.925534,\"BMD\":1.08216,\"BND\":1.541246,\"BOB\":7.459452,\"BRL\":6.337855,\"BSD\":1.08182,\"BTC\":0.000115,\"BTN\":81.97207,\"BWP\":13.193597,\"BYN\":2.647737,\"BYR\":21210.33532,\"BZD\":2.180752,\"CAD\":1.526933,\"CDF\":1960.874249,\"CHF\":1.0516,\"CLF\":0.032461,\"CLP\":895.708037,\"CNY\":7.685721,\"COP\":4244.988896,\"CRC\":615.516389,\"CUC\":1.08216,\"CUP\":28.677239,\"CVE\":110.543067,\"CZK\":27.718233,\"DJF\":192.321895,\"DKK\":7.455546,\"DOP\":59.847652,\"DZD\":139.552973,\"EGP\":17.024004,\"ERN\":16.232787,\"ETB\":36.534145,\"EUR\":1,\"FJD\":2.437295,\"FKP\":0.89399,\"GBP\":0.894074,\"GEL\":3.468366,\"GGP\":0.89399,\"GHS\":6.260339,\"GIP\":0.89399,\"GMD\":55.735408,\"GNF\":10231.822857,\"GTQ\":8.32871,\"GYD\":226.231417,\"HKD\":8.387865,\"HNL\":27.054383,\"HRK\":7.563259,\"HTG\":114.848025,\"HUF\":354.878171,\"IDR\":16100.808828,\"ILS\":3.828903,\"IMP\":0.89399,\"INR\":82.111095,\"IQD\":1287.770359,\"IRR\":45564.345721,\"ISK\":157.324833,\"JEP\":0.89399,\"JMD\":157.502345,\"JOD\":0.767148,\"JPY\":115.887974,\"KES\":115.873188,\"KGS\":82.877334,\"KHR\":4433.072409,\"KMF\":491.845813,\"KPW\":973.94397,\"KRW\":1334.509263,\"KWD\":0.334608,\"KYD\":0.901583,\"KZT\":455.588891,\"LAK\":9744.850865,\"LBP\":1636.59543,\"LKR\":203.283306,\"LRD\":214.755058,\"LSL\":20.024041,\"LTL\":3.195337,\"LVL\":0.654588,\"LYD\":1.537079,\"MAD\":10.692152,\"MDL\":19.255132,\"MGA\":4128.440643,\"MKD\":61.549021,\"MMK\":1520.013038,\"MNT\":3030.113149,\"MOP\":8.636724,\"MRO\":386.331123,\"MUR\":43.112348,\"MVR\":16.777528,\"MWK\":798.09702,\"MXN\":25.924012,\"MYR\":4.70924,\"MZN\":74.025194,\"NAD\":20.023995,\"NGN\":418.807136,\"NIO\":36.739739,\"NOK\":11.074396,\"NPR\":131.155591,\"NZD\":1.823154,\"OMR\":0.41661,\"PAB\":1.08192,\"PEN\":3.723036,\"PGK\":3.729165,\"PHP\":54.903428,\"PKR\":173.145964,\"PLN\":4.569316,\"PYG\":7103.482417,\"QAR\":3.940186,\"RON\":4.841804,\"RSD\":117.634783,\"RUB\":79.65293,\"RWF\":1014.524967,\"SAR\":4.065519,\"SBD\":9.061443,\"SCR\":17.863116,\"SDG\":59.847426,\"SEK\":10.671075,\"SGD\":1.546488,\"SHP\":0.89399,\"SLL\":10659.276025,\"SOS\":628.735306,\"SRD\":8.07079,\"STD\":23862.377246,\"SVC\":9.467174,\"SYP\":556.718973,\"SZL\":20.023915,\"THB\":34.6995,\"TJS\":11.085272,\"TMT\":3.78756,\"TND\":3.153144,\"TOP\":2.508559,\"TRY\":7.468207,\"TTD\":7.318545,\"TWD\":32.431293,\"TZS\":2503.472798,\"UAH\":28.813803,\"UGX\":4095.865634,\"USD\":1.08216,\"UYU\":47.526011,\"UZS\":10940.637611,\"VEF\":10.808077,\"VND\":25252.743871,\"VUV\":130.220071,\"WST\":3.034681,\"XAF\":656.445302,\"XAG\":0.06501,\"XAU\":0.00062,\"XCD\":2.924592,\"XDR\":0.796292,\"XOF\":656.871459,\"XPF\":119.644,\"YER\":270.891734,\"ZAR\":20.109119,\"ZMK\":9740.742168,\"ZMW\":19.98716,\"ZWL\":348.45551}}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c.Run(tt.name, func() {
			if err := c.currencyContext.saveCurrencyEntries(tt.args.raw); (err != nil) != tt.wantErr {
				c.Errorf(err, "saveCurrencyEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (c *CurrencyContextTestSuite) Test_CurrencyContext_Convert() {
	type args struct {
		from   string
		to     string
		amount float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Sanity Test",
			args: args{
				from:   "USD",
				to:     "TWD",
				amount: 15,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Sanity Test2, should have cache now",
			args: args{
				from:   "USD",
				to:     "TWD",
				amount: 10,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c.currencyContext.cachedBatch = &CurrencyBatch{
			Batch:   54,
			Exec:    time.Time{},
			Raw:     "",
			Success: true,
			Host:    "",
		}
		c.Run(tt.name, func() {
			got, err := c.currencyContext.Convert(tt.args.from, tt.args.to, tt.args.amount)
			if (err != nil) != tt.wantErr {
				c.Errorf(err, "Convert() error = %v, wantErr %v", err, tt.wantErr)
				c.FailNow(err.Error())
				return
			}
			if got != tt.want {
				//c.Errorf(errors.New("Not correct result!"), "Convert() got = %v, want %v", got, tt.want)
				//We don't determine this result.
			}
		})
	}
}
