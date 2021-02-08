package rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdata = []byte(`千字文 【作者】周兴嗣 【朝代】南北朝译文对照
天地玄黄，宇宙洪荒。日月盈昃，辰宿列张。寒来暑往，秋收冬藏。闰余成岁，律吕调阳。云腾致雨，露结为霜。金生丽水，玉出昆冈。剑号巨阙，珠称夜光。果珍李柰，菜重芥姜。
海咸河淡，鳞潜羽翔。龙师火帝，鸟官人皇。始制文字，乃服衣裳。推位让国，有虞陶唐。吊民伐罪，周发殷汤。坐朝问道，垂拱平章。爱育黎首，臣伏戎羌。遐迩一体，率宾归王。
鸣凤在竹，白驹食场。化被草木，赖及万方。盖此身发，四大五常。恭惟鞠养，岂敢毁伤。女慕贞洁，男效才良。知过必改，得能莫忘。罔谈彼短，靡恃己长。信使可覆，器欲难量。
墨悲丝染，诗赞羔羊。景行维贤，克念作圣。德建名立，形端表正。空谷传声，虚堂习听。祸因恶积，福缘善庆。尺璧非宝，寸阴是竞。资父事君，曰严与敬。孝当竭力，忠则尽命。
临深履薄，夙兴温凊。似兰斯馨，如松之盛。川流不息，渊澄取映。容止若思，言辞安定。笃初诚美，慎终宜令。荣业所基，籍甚无竟。学优登仕，摄职从政。存以甘棠，去而益咏。
乐殊贵贱，礼别尊卑。上和下睦，夫唱妇随。外受傅训，入奉母仪。诸姑伯叔，犹子比儿。孔怀兄弟，同气连枝。交友投分，切磨箴规。仁慈隐恻，造次弗离。节义廉退，颠沛匪亏。
性静情逸，心动神疲。守真志满，逐物意移。坚持雅操，好爵自縻。都邑华夏，东西二京。背邙面洛，浮渭据泾。宫殿盘郁，楼观飞惊。图写禽兽，画彩仙灵。丙舍旁启，甲帐对楹。
肆筵设席，鼓瑟吹笙。升阶纳陛，弁转疑星。右通广内，左达承明。既集坟典，亦聚群英。杜稿钟隶，漆书壁经。府罗将相，路侠槐卿。户封八县，家给千兵。高冠陪辇，驱毂振缨。
世禄侈富，车驾肥轻。策功茂实，勒碑刻铭。盘溪伊尹，佐时阿衡。奄宅曲阜，微旦孰营。桓公匡合，济弱扶倾。绮回汉惠，说感武丁。俊义密勿，多士实宁。晋楚更霸，赵魏困横。
假途灭虢，践土会盟。何遵约法，韩弊烦刑。起翦颇牧，用军最精。宣威沙漠，驰誉丹青。九州禹迹，百郡秦并。岳宗泰岱，禅主云亭。雁门紫塞，鸡田赤城。昆池碣石，钜野洞庭。
旷远绵邈，岩岫杳冥。治本于农，务兹稼穑。俶载南亩，我艺黍稷。税熟贡新，劝赏黜陟。孟轲敦素，史鱼秉直。庶几中庸，劳谦谨敕。聆音察理，鉴貌辨色。贻厥嘉猷，勉其祗植。
省躬讥诫，宠增抗极。殆辱近耻，林皋幸即。两疏见机，解组谁逼。索居闲处，沉默寂寥。求古寻论，散虑逍遥。欣奏累遣，戚谢欢招。渠荷的历，园莽抽条。枇杷晚翠，梧桐蚤凋。
陈根委翳，落叶飘摇。游鹍独运，凌摩绛霄。耽读玩市，寓目囊箱。易輶攸畏，属耳垣墙。具膳餐饭，适口充肠。饱饫烹宰，饥厌糟糠。亲戚故旧，老少异粮。妾御绩纺，侍巾帷房。
纨扇圆洁，银烛炜煌。昼眠夕寐，蓝笋象床。弦歌酒宴，接杯举觞。矫手顿足，悦豫且康。嫡后嗣续，祭祀烝尝。稽颡再拜，悚惧恐惶。笺牒简要，顾答审详。骸垢想浴，执热愿凉。
驴骡犊特，骇跃超骧。诛斩贼盗，捕获叛亡。布射僚丸，嵇琴阮啸。恬笔伦纸，钧巧任钓。释纷利俗，并皆佳妙。毛施淑姿，工颦妍笑。年矢每催，曦晖朗曜。璇玑悬斡，晦魄环照。
指薪修祜，永绥吉劭。矩步引领，俯仰廊庙。束带矜庄，徘徊瞻眺。孤陋寡闻，愚蒙等诮。谓语助者，焉哉乎也。`)

var testPubKeyPKCS1 = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDicXkhgbPXypYJlT51hxXGqXNI
Y4BtLFyYCtE3TEpfbA12TPHwGAG5A4aO7xORIbTJVFCmpxF/eFsVK4H1USC0NhtP
2rAvOfxGTu3Z1Fkp2hYt/e3OeAQbQze04QhvsIZzJKJV5AdcSYSTYGz5PGKw4TNr
kPvfsTAPdmMJcUAmAQIDAQAB
-----END PUBLIC KEY-----`)

var testPriKeyPKCS1 = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDicXkhgbPXypYJlT51hxXGqXNIY4BtLFyYCtE3TEpfbA12TPHw
GAG5A4aO7xORIbTJVFCmpxF/eFsVK4H1USC0NhtP2rAvOfxGTu3Z1Fkp2hYt/e3O
eAQbQze04QhvsIZzJKJV5AdcSYSTYGz5PGKw4TNrkPvfsTAPdmMJcUAmAQIDAQAB
An9pz9xqa9E95Mx3beXhxV3HlybnjJdWbIRYB7X7wQp/zF3+CbaEGrOzYJJf+BeM
mdOAwIVvWmIyzmrBWbNAjshYjCluNVXD+tB2chcMPFXF9mimjy2+EC/hkDVBWjcA
SRpY3/NDZwTqGEJz4NQLwVWAy1UJobjfGhhgfQL2eCxhAkEA/KJlWCElZJzmP3Dn
lSKzgk5hlafwnE2y+6Hg0Tr/MsMILf5IIw0clKE1L54n3ceJnH7nRBkVpccnGjCH
FjJGcQJBAOV1wMo3WnsNpIUTkHJJ+IPL1hsP/OfmrUtoNMSyvDGR6c6udVM7YF2T
0XT0bjcBUlROJ9JZSTeedBXQtTjKQJECQC0MMBIE5wwHxi6tzT2UkHm9zDzJU2gr
mqyv8sycPoEosb6xxt8pKV1/WWKCSw2K1QjowAguiOOknV5YJN5aXKECQQCt7Go8
HGbVvMqGIAUty8m1xGw+SQkOkbeq34qXyU6CWDIbefruIqRxaZirCJb91F+eDTt7
4jdwFAezfWXPbOYxAkEAwpKqeQCko4J8HBJrAjIdEx+/K4HJgONBEU1FnujQ0QwA
HKPEnWKZQj7e9fcCQ4ped1bkUm/vOI6HSeNtg6A/VQ==
-----END RSA PRIVATE KEY-----`)

var testPubKeyPKCS8 = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy0RfiBzH5SLKg85yFVGm
qO5u4AkedWB+VvNB+yoRjvwqFZumG5xdVSKwOlPIAWY4QsvwuZ0DV7BIIzu4Rova
ABzBzyyLswxuoihrEwLU34azsCVJ8L8309SFMhfgnAQo8rSCKuapL9TU2epyGTyj
frEv1zlLKl/eghVFlRI8rIGk5SuCRYMRu4NXbF/v5M1URyu+D4GdVd2eH3WB3ZXs
dM5GcO29sZ7C3LLoLX7BooAYg2TGuVWnDkNeUOn5mNDi3dYWC5kvLYF7Jv1wf1lv
osLj0KmvBBziqg22jujfmcf2drAUS9APB0scJwETzfPFRbSoIgJG05lq8ZrH5H7c
9wIDAQAB
-----END PUBLIC KEY-----`)

var testPriKeyPKCS8 = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDLRF+IHMflIsqD
znIVUaao7m7gCR51YH5W80H7KhGO/CoVm6YbnF1VIrA6U8gBZjhCy/C5nQNXsEgj
O7hGi9oAHMHPLIuzDG6iKGsTAtTfhrOwJUnwvzfT1IUyF+CcBCjytIIq5qkv1NTZ
6nIZPKN+sS/XOUsqX96CFUWVEjysgaTlK4JFgxG7g1dsX+/kzVRHK74PgZ1V3Z4f
dYHdlex0zkZw7b2xnsLcsugtfsGigBiDZMa5VacOQ15Q6fmY0OLd1hYLmS8tgXsm
/XB/WW+iwuPQqa8EHOKqDbaO6N+Zx/Z2sBRL0A8HSxwnARPN88VFtKgiAkbTmWrx
msfkftz3AgMBAAECggEBALFU8h8WNNJTWuhIuECfHk2szfAIJJl0WGRVDtAnMHBU
0AZP50eZT9eRTrtkEk4eNqEXuOjF8X+T3CtY03hAVCza6b5rd2j6RJ6VsmxAgSSN
KMBYl7w/pl3Yv2cna53RB0ROdG0nFJ0VCakfnvEjdON3g2w5oLjUMJO+KRlXcrVv
Szru6JJbz5xlVmLRAg4afEZWvXUjARmHjakgankxUof3e+hIPVcfVdHxclCEIAJz
FIPM74LLR0kT7Rh1W1eN53DMD5AAFMSe7xk0U5qvYK0ARiZsfOssJOEoGRf0iJPx
gLgTNNXPeLVkWGwpI2oWC9jujZwiFSc76EcnXcTsOMkCgYEA5ZzGi8w5cF/7TxCI
HLQzh+zWnv4orQwLM7xtF15lBuDCRcQPotN4nkq/f69v/DS/WPOJmxp/KvTxkNpN
lB9Hm+bznTNozVAPL/k2arrPQagr+Rr7AyO9fsIHWK8x3estlnxp5OgGCHpGhEhz
qd8aFSIKwq9/0HW0HesGBsRuab0CgYEA4qCD3BB8aNioyCxibQyN/CaYWytnf5MH
wmpick6shpiglX1o/TmoH9hagWjmTkxcBeUOXVE2ch8ZWYfY4jIFU2roSonBCrLh
YHNh6I582TkRinskD3rDzVF/HBzll6YPQ/tZXfboq++qYuFsp2qyBt0egDsTRrmW
+6NknORlusMCgYEAwsCF+x8+fN08SCSKfmYt6xVsOLD0iUpU2g3xLcgHwpyyyv/w
Dzh07zYRVVjVkEKhJe5zAdaolCPsHOO8t20MjOSILwbi1noYV6V2jXJjxpnAqmpD
C7ety60BFCyNDGCkayadnuVZ8Kjd1OerCyNLeS9FlznKHGMpYdLtaqID0LUCgYA2
HT//L0yVuI5s5fRGt8W7nPeqZW3cT559tOt3AgQ+S3mk2IJWXQshN4c8+XBs59zd
Z3mLnNXUYEqsTzzhnjIZXiDDk6stw9L/Ne3+GvAC6paeq5LLw3O3tisU6m2ETZm9
kOog/tFGJP9ZhxxryZVjAC/FTNXogG5l/fkLYZpNAwKBgDKG6rjVUJf5FxpEcvV7
JJBmRMP6XHycESytqddAb6WZ3uWhA7jYH4Lmzasx7UvzULDr4NbL/dvyuc4psgWO
Ud6dDG+l0otkuKxVEisl4yHyhMDFcf3YVn+7zz0NlZFflC3prph3pZ+AYACZ5gJ+
9l01SZfIKVdosE3WTyynjWNi
-----END PRIVATE KEY-----`)

func TestXRsa(t *testing.T) {
	opts := DefaultOptions()

	opts.KeyType = PKCS1
	testCreateKeys(t, opts)

	opts.KeyType = PKCS8
	testCreateKeys(t, opts)

	testPKCS1(t, testPubKeyPKCS1, testPriKeyPKCS1)
	testPKCS8(t, testPubKeyPKCS8, testPriKeyPKCS8)
}

func testPKCS1(t *testing.T, pubKey, priKey []byte) {
	opts := DefaultOptions()
	opts.KeyType = PKCS1
	testXRsa(t, pubKey, priKey, opts)
}

func testPKCS8(t *testing.T, pubKey, priKey []byte) {
	opts := DefaultOptions()
	opts.KeyType = PKCS8
	testXRsa(t, pubKey, priKey, opts)
}

func testCreateKeys(t *testing.T, opts *Options) {
	pubKey, priKey, err := CreateKeys(opts)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, pubKey)
	assert.NotEmpty(t, priKey)

	testXRsa(t, pubKey, priKey, opts)

	//t.Logf("\n%s", pubKey)
	//t.Logf("\n%s", priKey)
}

func testXRsa(t *testing.T, pubKey, priKey []byte, opts *Options) {
	xrsa, err := NewXRsa(pubKey, priKey, opts)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, xrsa)
	assert.NotNil(t, xrsa.opts)
	assert.NotNil(t, xrsa.publicKey)
	assert.NotNil(t, xrsa.privateKey)

	testPKCS1v15(t, xrsa)
	testOAEP(t, xrsa)
	testPSS(t, xrsa)
}

func testPKCS1v15(t *testing.T, xrsa *XRsa) {
	xrsa.opts.Type = PKCS1v15
	testEncryptAndDecrypt(t, xrsa)
	testSignAndVerify(t, xrsa)
}

func testOAEP(t *testing.T, xrsa *XRsa) {
	xrsa.opts.Type = OAEP
	testEncryptAndDecrypt(t, xrsa)
}

func testPSS(t *testing.T, xrsa *XRsa) {
	xrsa.opts.Type = PSS
	testSignAndVerify(t, xrsa)
}

func testEncryptAndDecrypt(t *testing.T, xrsa *XRsa) {
	encode, err := xrsa.Encrypt(testdata)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("encode:", encode.Base64())

	decode, err := xrsa.Decrypt(encode)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("decode:", string(decode))

	assert.Equal(t, decode.Bytes(), testdata)
}

func testSignAndVerify(t *testing.T, xrsa *XRsa) {
	signcode, err := xrsa.Sign(testdata)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("Sign:", signcode.Base64())

	if err := xrsa.Verify(testdata, signcode); err != nil {
		t.Fatal(err)
	}
	//t.Log("Verify:", err)

	assert.Nil(t, err)
}
