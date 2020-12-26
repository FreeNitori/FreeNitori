package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"strings"
)

type CharacterInfo struct {
	Color        int
	FriendlyName string
	SearchString string
}

func CharacterList() []CharacterInfo {
	return []CharacterInfo{
		{Color: 0xb50404, FriendlyName: "Reimu", SearchString: "Hakurei_Reimu"},
		{Color: 0x42d4f4, FriendlyName: "Ran", SearchString: "Yakumo_Ran"},
		{Color: 0xff8ade, FriendlyName: "Mystia", SearchString: "Mystia_Lorelei"},
		{Color: 0x42d4f4, FriendlyName: "Miyoi", SearchString: "Okunoda_Miyoi"},
		{Color: 0x42d4f4, FriendlyName: "Mamizou", SearchString: "Futatsuiwa_Mamizou"},
		{Color: 0x42d4f4, FriendlyName: "Eternity", SearchString: "Eternity_Larva"},
		{Color: 0x0, FriendlyName: "Nue", SearchString: "Houjuu_Nue"},
		{Color: 0x42d4f4, FriendlyName: "Konngara", SearchString: "Konngara"},
		{Color: 0x42d4f4, FriendlyName: "Rinnosuke", SearchString: "Morichika_Rinnosuke"},
		{Color: 0x42d4f4, FriendlyName: "Yumemi", SearchString: "Okazaki_Yumemi"},
		{Color: 0x42d4f4, FriendlyName: "Lily", SearchString: "Lily_White"},
		{Color: 0xdf9041, FriendlyName: "Kutaka", SearchString: "Niwatari_Kutaka"},
		{Color: 0x42d4f4, FriendlyName: "Medicine", SearchString: "Medicine_Melancholy"},
		{Color: 0xf5e942, FriendlyName: "Marisa", SearchString: "Kirisame_Marisa"},
		{Color: 0xd25859, FriendlyName: "Raiko", SearchString: "Horikawa_Raiko"},
		{Color: 0x4e7764, FriendlyName: "Mai", SearchString: "Teireida_Mai"},
		{Color: 0xa84384, FriendlyName: "Yorihime", SearchString: "Watatsuki_No_Yorihime"},
		{Color: 0x42d4f4, FriendlyName: "Nemuno", SearchString: "Sakata_Nemuno"},
		{Color: 0x42d4f4, FriendlyName: "Suguri", SearchString: "Suguri_(character)"},
		{Color: 0xc7c7c7, FriendlyName: "Sakuya", SearchString: "Izayoi_Sakuya"},
		{Color: 0x42d4f4, FriendlyName: "Wan", SearchString: "Inubashiri_Momiji"},
		{Color: 0x42d4f4, FriendlyName: "Kisume", SearchString: "Kisume"},
		{Color: 0x42d4f4, FriendlyName: "Star", SearchString: "Star_Sapphire"},
		{Color: 0x42d4f4, FriendlyName: "Mima", SearchString: "Mima"},
		{Color: 0x9917, FriendlyName: "Okuu", SearchString: "Reiuji_Utsuho"},
		{Color: 0xb2daef, FriendlyName: "Nitori", SearchString: "Kawashiro_Nitori"},
		{Color: 0xaa6ad3, FriendlyName: "Sumireko", SearchString: "Usami_Sumireko"},
		{Color: 0x42d4f4, FriendlyName: "Sukuna", SearchString: "Sukuna_Shinmyoumaru"},
		{Color: 0xf5da42, FriendlyName: "Ex_Rumia", SearchString: "Ex-Rumia"},
		{Color: 0x42d4f4, FriendlyName: "Kokoro", SearchString: "Hata_No_Kokoro"},
		{Color: 0x42d4f4, FriendlyName: "Futo", SearchString: "Mononobe_No_Futo"},
		{Color: 0x42d4f4, FriendlyName: "Mayumi", SearchString: "Joutouguu_Mayumi"},
		{Color: 0x42d4f4, FriendlyName: "Renko", SearchString: "Usami_Renko"},
		{Color: 0x42d4f4, FriendlyName: "Wriggle", SearchString: "Wriggle_Nightbug"},
		{Color: 0x14a625, FriendlyName: "Kokuu", SearchString: "Kokuu_Haruto"},
		{Color: 0x42d4f4, FriendlyName: "Tenshi", SearchString: "Hinanawi_Tenshi"},
		{Color: 0x79eb50, FriendlyName: "Youmu", SearchString: "Konpaku_Youmu"},
		{Color: 0x42d4f4, FriendlyName: "Letty", SearchString: "Letty_Whiterock"},
		{Color: 0x42d4f4, FriendlyName: "Minoriko", SearchString: "Aki_Minoriko"},
		{Color: 0xe2a81e, FriendlyName: "Ringo", SearchString: "Ringo_Touhou"},
		{Color: 0x24b343, FriendlyName: "Sanae", SearchString: "Kochiya_Sanae"},
		{Color: 0x940f0f, FriendlyName: "Hecatia", SearchString: "Hecatia_Lapislazuli"},
		{Color: 0x42d4f4, FriendlyName: "Merlin", SearchString: "Merlin_Prismriver"},
		{Color: 0x42d4f4, FriendlyName: "Luna", SearchString: "Luna_Child"},
		{Color: 0x42d4f4, FriendlyName: "Chiyuri", SearchString: "Kitashirakawa_Chiyuri"},
		{Color: 0xa700f5, FriendlyName: "Satori", SearchString: "Komeiji_Satori"},
		{Color: 0xfd8cff, FriendlyName: "Remilia", SearchString: "Remilia_Scarlet"},
		{Color: 0x42d4f4, FriendlyName: "Hina", SearchString: "Kagiyama_Hina"},
		{Color: 0x42d4f4, FriendlyName: "Eirin", SearchString: "Yagokoro_Eirin"},
		{Color: 0xe58a53, FriendlyName: "Aya", SearchString: "Shameimaru_Aya"},
		{Color: 0x42d4f4, FriendlyName: "Kagerou", SearchString: "Imaizumi_Kagerou"},
		{Color: 0x42d4f4, FriendlyName: "Lunasa", SearchString: "Lunasa_Prismriver"},
		{Color: 0xe5ff, FriendlyName: "Cirno", SearchString: "Cirno"},
		{Color: 0x42d4f4, FriendlyName: "Sariel", SearchString: "Sariel"},
		{Color: 0xf50000, FriendlyName: "Mokou", SearchString: "Fujiwara_No_Mokou"},
		{Color: 0x42d4f4, FriendlyName: "Suika", SearchString: "Ibuki_Suika"},
		{Color: 0x42d4f4, FriendlyName: "Sekibanki", SearchString: "Sekibanki"},
		{Color: 0x42d4f4, FriendlyName: "Urumi", SearchString: "Ushizaki_Urumi"},
		{Color: 0xb50404, FriendlyName: "Flandre", SearchString: "Flandre_Scarlet"},
		{Color: 0x42d4f4, FriendlyName: "Aunn", SearchString: "Komano_Aun"},
		{Color: 0xd25859, FriendlyName: "Komachi", SearchString: "Onozuka_Komachi"},
		{Color: 0xaeb4c6, FriendlyName: "Seija", SearchString: "Kijin_Seija"},
		{Color: 0xb50480, FriendlyName: "Hatate", SearchString: "Himekaidou_Hatate"},
		{Color: 0x42d4f4, FriendlyName: "Momiji", SearchString: "Inubashiri_Momiji"},
		{Color: 0x42d4f4, FriendlyName: "Parsee", SearchString: "Mizuhashi_Parsee"},
		{Color: 0xef61ff, FriendlyName: "Kaguya", SearchString: "Houraisan_Kaguya"},
		{Color: 0x990000, FriendlyName: "Koakuma", SearchString: "Koakuma"},
		{Color: 0xe262b0, FriendlyName: "Satono", SearchString: "Nishida_Satono"},
		{Color: 0x583b80, FriendlyName: "Toyohime", SearchString: "Watatsuki_No_Toyohime"},
		{Color: 0xf94aff, FriendlyName: "Reisen", SearchString: "Reisen_Udongein_Inaba"},
		{Color: 0x42d4f4, FriendlyName: "Kyouko", SearchString: "Kasodani_Kyouko"},
		{Color: 0x42d4f4, FriendlyName: "Yoshika", SearchString: "Miyako_Yoshika"},
		{Color: 0x42d4f4, FriendlyName: "Seiga", SearchString: "Kaku_Seiga"},
		{Color: 0x42d4f4, FriendlyName: "Miko", SearchString: "Toyosatomimi_No_Miko"},
		{Color: 0x2291ba, FriendlyName: "Rei", SearchString: "Reisen"},
		{Color: 0x42d4f4, FriendlyName: "Kogasa", SearchString: "Tatara_Kogasa"},
		{Color: 0x24b343, FriendlyName: "Yuuka", SearchString: "Kazami_Yuuka"},
		{Color: 0x42d4f4, FriendlyName: "Saki", SearchString: "Kurokoma_Saki"},
		{Color: 0x42d4f4, FriendlyName: "Shinki", SearchString: "Shinki"},
		{Color: 0x574b8c, FriendlyName: "Keine", SearchString: "Kamishirasawa_Keine"},
		{Color: 0x42d4f4, FriendlyName: "Yachie", SearchString: "Kicchou_Yachie"},
		{Color: 0xb4449c, FriendlyName: "Sagume", SearchString: "Kishin_Sagume"},
		{Color: 0x42d4f4, FriendlyName: "Yamame", SearchString: "Kurodani_Yamame"},
		{Color: 0x42d4f4, FriendlyName: "Chen", SearchString: "Chen"},
		{Color: 0x4b548, FriendlyName: "Meiling", SearchString: "Hong_Meiling"},
		{Color: 0x42d4f4, FriendlyName: "Iku", SearchString: "Nagae_Iku"},
		{Color: 0xc646e0, FriendlyName: "Patchouli", SearchString: "Patchouli_Knowledge"},
		{Color: 0x42d4f4, FriendlyName: "Nazrin", SearchString: "Nazrin"},
		{Color: 0xcc7c9c, FriendlyName: "Tewi", SearchString: "Inaba_Tewi"},
		{Color: 0x42d4f4, FriendlyName: "Eika", SearchString: "Ebisu_Eika"},
		{Color: 0x42d4f4, FriendlyName: "Shou", SearchString: "Toramaru_Shou"},
		{Color: 0x42d4f4, FriendlyName: "Wakasagihime", SearchString: "Wakasagihime"},
		{Color: 0x42d4f4, FriendlyName: "Alice", SearchString: "Alice_Margatroid"},
		{Color: 0x42d4f4, FriendlyName: "Genjii", SearchString: "Genjii"},
		{Color: 0xaa4fa0, FriendlyName: "Joon", SearchString: "Yorigami_Joon"},
		{Color: 0x977cac, FriendlyName: "Kanako", SearchString: "Yasaka_Kanako"},
		{Color: 0x42d4f4, FriendlyName: "Lyrica", SearchString: "Lyrica_Prismriver"},
		{Color: 0x42d4f4, FriendlyName: "Keiki", SearchString: "Haniyasushin_Keiki"},
		{Color: 0x42d4f4, FriendlyName: "Orin", SearchString: "Kaenbyou_Rin"},
		{Color: 0x42d4f4, FriendlyName: "Sunny", SearchString: "Sunny_Milk"},
		{Color: 0x4b548, FriendlyName: "Daiyousei", SearchString: "Daiyousei"},
		{Color: 0x42d4f4, FriendlyName: "Maribel", SearchString: "Maribel_Hearn"},
		{Color: 0x5b0082, FriendlyName: "Byakuren", SearchString: "Hijiri_Byakuren"},
		{Color: 0x5b9c66, FriendlyName: "Eiki", SearchString: "Shiki_Eiki"},
		{Color: 0xe5ff, FriendlyName: "Vert", SearchString: "Cirno"},
		{Color: 0xfb959e, FriendlyName: "Kasen", SearchString: "Ibaraki_Kasen"},
		{Color: 0xff40d9, FriendlyName: "Yuyuko", SearchString: "Saigyouji_Yuyuko"},
		{Color: 0x42d4f4, FriendlyName: "Clownpiece", SearchString: "Clownpiece"},
		{Color: 0x42d4f4, FriendlyName: "Doremy", SearchString: "Doremy_Sweet"},
		{Color: 0x42d4f4, FriendlyName: "Yuugi", SearchString: "Hoshiguma_Yuugi"},
		{Color: 0x42d4f4, FriendlyName: "Yukari", SearchString: "Yakumo_Yukari"},
		{Color: 0x42d4f4, FriendlyName: "Youka", SearchString: "Kazami_Youka"},
		{Color: 0xe69454, FriendlyName: "Okina", SearchString: "Matara_Okina"},
		{Color: 0x62f500, FriendlyName: "Koishi", SearchString: "Komeiji_Koishi"},
		{Color: 0x42d4f4, FriendlyName: "Suwako", SearchString: "Moriya_Suwako"},
		{Color: 0x42d4f4, FriendlyName: "Shizuha", SearchString: "Aki_Shizuha"},
		{Color: 0x42d4f4, FriendlyName: "Ichirin", SearchString: "Kumoi_Ichirin"},
		{Color: 0xf5da42, FriendlyName: "Rumia", SearchString: "Rumia"},
		{Color: 0x42d4f4, FriendlyName: "Narumi", SearchString: "Yatadera_Narumi"},
		{Color: 0x42d4f4, FriendlyName: "Kosuzu", SearchString: "Motoori_Kosuzu"},
		{Color: 0x6b87bd, FriendlyName: "Seiran", SearchString: "Seiran_(touhou)"},
		{Color: 0xfbd55a, FriendlyName: "Junko", SearchString: "Junko_(touhou)"},
		{Color: 0x42d4f4, FriendlyName: "Murasa", SearchString: "Murasa_Minamitsu"},
		{Color: 0x42d4f4, FriendlyName: "Akyuu", SearchString: "Hieda_No_Akyuu"},
		{Color: 0x48cb5, FriendlyName: "Shion", SearchString: "Yorigami_Shion"},
		{Color: 0x42d4f4, FriendlyName: "Elly", SearchString: "Elly"},
	}
}

func init() {
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "touhou",
		AliasPatterns: []string{"t", "th"},
		Description:   "Finds picture of requested character.",
		Category:      multiplexer.MediaCategory,
		Handler:       touhou,
	})
}

func touhou(context *multiplexer.Context) {
	if len(context.Fields) < 2 {
		context.SendMessage("Please specify a character.")
		return
	}
	name := strings.ToLower(context.Fields[1])
	for _, character := range CharacterList() {
		if name == strings.ToLower(character.FriendlyName) {
			// TODO: Implement image fetch
			embed := embedutil.NewEmbed(character.FriendlyName, "")
			embed.Color = character.Color
			embed.SetImage("")
			embed.SetFooter("Source URL: " + "")
			context.SendEmbed(embed)
			return
		}
	}
	context.SendMessage("Your request did not match any character, please try again.")
}
