package adapters

import (
	"context"
	"fmt"
	"log"
	"strings"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) ActivateBlockedRanges(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs := []string{
		"1c1yJLIt3bxbJc8ykv12RkUSjngGOI-Z_8qoaTygwhqo", "18Pc0NTmkr9kPExWtCFmlWtRggySMzRl5PGeAbaAprSo", "1JkK7renOvJ-mDTojLP_k1bjsu8sTRpr9MyERZHq1Jt0", "1UvEbSrUtN-EXY7Qc4gAXWtmzQyestRlw2H97DyCNpjw",
		"1znzdQsKIHqVEx9J_9n8iOzLs67OyPP8uCDXzljmDdT0", "1ZUPZSzElxz50owyBRBunxpi2AhoBwKmWEa3xWfHWZGg", "1Xf8bvGBt63w6TS2G3k3EN-R7c8ujSipcUSdy2fHU2MY", "1cesfjKCRezQoGSzX_fj_y8H2rO_5BQWaPG0JQ9h7hp8",
		"1rw1rqoTWHGEjTNkQt0iL1OZGyTtXsVarAPV3J16uwQo", "1nVTgSJCalLDfR2prFgdYV2IydMXcGiY4o9Fv3aCaSYg", "1r7EzdxwCfcCGD9MQzekx0yEnghofppMDQQtXXj_KSzs", "162EpE6JVTkkY4pE9XuhjTuNuC_bz6LHvTN7I5ReXQnQ",
		"18K_1Ls2UthlndfHK4YkQQtvzsvNEKAnhYw70qkN_xog", "1HsEdbx-BxLn-k5fAiOmY5dl8xazDbE4fe0E6dIx2mpA", "1eLCHxSZW-YZwgy_erhAd-kodEfLoOmHOzgHj1CHgEus", "14btaJqWcSWqV5PSaXED6kkSDXibvurGialcwkyZAGCU",
		"1WfMOr_QQhxZm6lpDSjIee-p99euhIjCBKkEbVTw4RbY", "1X1C3YS9BfPAmldenCs4ij3qPVZo_fYt9LwcM9sDKLrc", "1Ta1y1660bzM8YpXiFly_cVW9AGQgQCD25TCew_vkfyU", "1lYpH3QpJOAio9jD7HQ1GA5cIgBRYNuIX-TH2gDmGlcM",
		"1TZwHYSI3975St4bld0yG2fngEmsXoLH0afOLlzAciQs", "1VazpGDgp7IQsXuoGzDuWfG20eZuXTmluClb4TiCbWUA", "1qzB8mLsvc0dl9yC3gTskpmnubQwn__1ipT89W_s32YA", "10ZwKF9DB0bD3qEjJdKjGjM5ZIo3uc1NDYQBmYOPwlhc",
		"1AisI_Xz0PgE0OuuMWFF2WPK73RJzPRdG3BFRi7bHCJk", "1jHhf9zF-nPBh-rnQSpoXk1g2n830X6aXwbjW5k4wqTU", "1tEaZczQMwLZbcY6j9p4JzrRl0vTRh3MUNWGWUbj7Uv0", "1qihHg_NZ6G70140KMfi3FwDLBi5B2rQm1SfGdM5wPgs",
		"1sAPLUcbNqZtw7W1eoOu82oxDcqnM5uAxIx_-YHIopEU", "1Jp4sTkvvEumPFuZA7fpa1Z1U3R9i2bI9V3-z60duNWg", "1pfUsE5MWHMMocAqKHTo3wYHYlgIoluoN5DZsAs9k9vk", "12r7wG6tT5prPVuHuk92cax9kICFMRK1t165clZv18W4",
		"1kEnY1tcV3mEjfBXAFOypqTvwmbIYgnECzQVRZkAgDso", "1fvfFht0bNsKkHHQztWQKNq-AjEVu8Sj5HVdn5irwfnM", "1V0p80kryN9FpW6Hpw4f6RfqXLsBPaF9dsXCVV2Ljm4U", "1oGUL6JYfoMEo6DJNukyl6neFVvYcbZwVt6GfDhrcwjs",
		"1vnac2RXonPGLvaM697WfeZXjTjBXY0SNx8jm14T7jE0", "1y-N0SKgI5L-qJCQ9k4EhJN-PAQ61uJSOjhCtUKrH94g", "1k9Ray-9O9cvZullkBq0rLR0r8xw_mwlu6AVP7ieOiZM", "1zzu4rB20EjQbQXzhYgB1iGso59H-szZ8l8i4IENOgm8",
		"1Kc6M4K6R-PGO0tiki4LpHDY20S4TvjODDfG-39L8_vI", "1QzIHgTSS6ZFJGo_cN_JRCnJVi1m0xjxFoYE_8ux6QaI", "1xyulTOre1Tik4bbtTQAk-Jk36y-JeRSGhf8xX0Zs1-k", "1p50HLUkBpACGIsiBCGKKUpTr21tXHShYSuYQlo-Hs-k",
		"1XezdKLE4WBgmFGMzYwhIxbZEWjRYekOhdeeMDBS8s9U", "19FGj2bC1UtqcnZy3VCwIRaAmbqEHkVQMhXIoy1RonVY", "1vbnYblfnyJ8OJ7bzywTeVXDZQoiBXHjD1WM5WXfFKPw", "1GGrPEs5gbl1r-MErgJM6Sx5FCVrSp6U5rmD8ILJk_PE",
		"1g2Rl-4tCfhtn7_HpXrPkLiC60qRYISuBr9XFWOUvWYc", "115px4QhFpsSRW40rtHkLK5IxerzoSB7fCgpbYnN7bcs", "1fTLQrAmiHHADHnyc2xRREFjZlaaHLs4OQ8BY38axJgI", "1Qddmwvv0CHmgVxP3N3KdLNvxDMA8yEypFu7oRlS8N_0",
		"1MIxbo89FzVwy14VLhK7N0OGxofoo9_rIUNVeDI9F9NI", "1n0WapvYle1aSCTCg5C3LH9snIxhJL5zvm4BvwcOkrnU", "12FKiNdd8cudtAEdNCsoNIE0xlCoYNqj7tnqUAv9lWME", "1hyjli8bQB_hDervJ796blv7PrRVE5DLsX6E_QRmSp0U",
		"1k1qv_oX7EicScjWRvYe3-eKjLMmkyyz3-g3D83VWGe4", "1mQTnR9Itbyf6VX8DGj0zvR4F9QXOXOVVkRp9ODhIngE", "1hxAvB6t5wYf7GZYlBDhUeg_zZQ0RPTGR9jZblMp7maY", "1-yaMCDLWfrZPciiJlBzKOcAhQZlS9TEP7w8SnhtSC8U",
		"1dnxGao23OFWdT0v_ry5AcGn9Rrk7DQP62KO6Kg5Xymw", "1KMyJp-JVjkNPJ1vNzkg18dNusvYRj6T4XZTTAJ6Eiuo", "1t8W2lQEz64jDta5XGK3LSb-c36GOmZ-arMCXx6fBzRU", "1J0dNd4yxxa-DWmUIWUNDdnJjVtYIYsRpL98JZdZLNUM",
		"1pa0rVaX3sK4pZSu9o9RCrQwqHZpDXDjrlgztjjLAOfg", "1X67Olsph3fEqxsDpwETWj5UvzWxpLuUlrrC1mww1HZI", "1OdfMSikfzh6NvAq2rs1bFVWgs9PWVtUBCN4gCvZgW_I", "15ez9Hs8SF0RMhcAPRnzs7ptcU0j0_fYesJLOvZA3MNo",
		"1baKhMgeoNEH5TEe4H0fE5gnJPCgn_8wcJw_lCxmP8CU", "1mPvZCp4RrtMjGQprmSjNCmr9GiklrFwVJOfZNWYiDaI", "1ZICQd7hfJdVtiVsuSAew2dcPy9NXzf0zi0Xx1ulJ8SA", "119ygHqSgwjPwOekTQm2RayiCPQ00LnUtZQXjUHgQEEA",
		"1yqqhxZHfujAkUrnRW11MetPmBpzp-bMaN7Okk6y60-U", "1qGoWvPjN5w435Zm_lkt1GbcUM2VZZ5w_b5Gl_6i_ou8", "1MdqjufVjUA6pJomm1FjZ5g6s4RHoxSlbYX7N0vzswfo", "1dK6trumK_tf1Amy3jEBnHPD9Oabmv87J9arsPIgEvtY",
		"1k5yP6XW11nUVnS-LioJE18FbxH1S_GvhNJ8j9DSdsaI", "1bchRvC25IvRL12EZ_a0Y5ZIXtBm3gUIrxtlYr1Pqzlw", "1Y2nfJqPrPA7J-7d1cxHfRTE5wdkEQ2wg4oK-uH-kVLY", "1nHk5OJbeIGs2l4MHJcPI7fMYtNFddxSqbEuwBDHtwJ0",
		"1IPQx2bCbJ6TckvqI1LVgwUugb4Vj4HVLW53QT82fdTg", "19OoAnsYIInnu3WdByKb0oEQ7zC1i2BvqTSAbdlnp43Y", "1xfKRzSYLtN0GeLyIoPbu5xaI8foFpFbEZWIdX68_VkI", "1apEzPvVSgxPZshTz1TGF2QNgosZE2C1jBw0j_wjO3lQ",
		"1JHavaD9ZrrLDYZotJVd9N1B9eAMIa9XzGnZVay2hmJU", "10mytod7TsGLZJ6vASAid3dnpa6OFoZ4ZujQW5NpqUdc", "1wE4JGBT7uaZvt4J4VOucv9pnDsUxw0OwQK__IgQcQ0k", "1sL0iNHrZ4BKnXuSWAkQImz0hBwQPLLl7xda8aDHS888",
		"14ipQ1hxSoyTpkHQGtg4DK-KMCJJV9sDbR7rOeLbkJyI", "1wJvOZB215ijgALSi1NlMZRnQnWwkmQETMlh3X3tfxVg", "1Rt047s1JCdRqnlaBZ9af5jobcBK-V9QSZjqGyt3cv6E", "1cjMRYw1v5MxRoAuldhEXGjAiNhQr36VmttjhGSECjs0",
		"1CYF8AOL2dI2yY-JhqxZzEnT50BK5Op_i228By6NsRWI", "1uzkWstoM-sDWNzaNaaQecOg7DjpqEs3A2LHRl6R3Qv8", "1d5Ndtdptd_Jcstm3jEeVsx31h61qP0EghJygL2EaFzk", "1VHyaTWbvltfshMoPemkTjZEuYdbl95Yd_0-tont6c7c",
		"13BNidTuWqh4WvPsv-hFZhLtZWrITqw1T2lbG3UjRmQQ", "1IE8e7FM7CUD_gskwbGqyka43KFWfpSAvPX2vllUQmDI", "1ZPft2_E-4eXgdtByOc2k1YP2lf-7NWzPKjRW0-ZVCEI", "10TA83SY50q13Wp7Ip0-Y8j4ugk3ZEhxFLy8rbra3Eq4",
		"1BurOnFFSOCFcMFSpjnZzzBJVNlhPkszeH43erAv1Ajg", "1245aNoMJBqb4wTVgXe3Qo0EuD2f1vTxOIrLctCetrsQ", "1rrD7WSY9pSr4lsjpnJ_Qo1dRxkqyy33coVQyXI5cvro", "1cD-eh5fy2wY6zdSzN9bLm-rNJ8IhVTHJcqseY9b8eR0",
		"1uBG1zftgHn_uhYmXPgk0ClyAn7zmHA8PCeIXYe3JxaE", "190U6SzY8nd1gpRcH5nhlkwmyq5A7W0rLABuBIZhOVyI", "1KQ0Ye4_BSnPdjD5mODiTFOo2Kt7oZHzPBSJ5TEA9y9c", "1OznazxS69I6nwxOgYWi3R7-vTIngt-2Thymy8J_XEgg",
		"1ZgGtPd3kjD6N_Q-Cmks4p_wxhHVTKvge2hzxP_i0riA", "1l7ACAWjJUY9NALboiOxu6NenS22dp1F3bIH0CmkW2dM", "1V_u3mLrxNMGHLGrw-6DwbG4NBXEdYiBnUH_VaxnpT7M", "1Gs3l1o5cqN49p8-2t7iZaYQ0CV3zp_5P5C57DliGoEA",
		"1y_B0iBNYPjauPkbqU7M64mwR7KB6_nIvYHYYp51jn04", "1lZ4c-8xQc1Cz9snqFghloxx0qKNXq9nrLs9LucNIkGs", "1JLVHgocleGmsHX9y_Q0e675M5Wcw3zvkT6a0nHoxlp4", "1LstA-TDRHjJ9BL9_xmR1vBM4LUCWr3qj5PF_TfUn7ko",
		"16EEyi7Dv8qI6fgdWZUyTN70MafiyecWFNbmxVXJJBXs", "1n0oxcNlKey1v-0ymRdRfSn6Cx7IrfYbU25nRPUEUmtQ", "1nDWIR0splwvay39TAQklprOHY02XettoSMz3NE174m0", "12-DvqvJSohPsCB3iFDprXGcCRYqL6F7hU5oPBOB7hDs",
		"1XCeSvM2JAUXaVH4nTbMYPfPjAMb3Y_JCd4qH54M-r3k", "1SSe-zjdSITh3QNAMXiRoDth12feOY2qipLXB_DzFQxI", "1CoVi36PW6H3Xqu4XZV-ONvIOlP9Dhm4DY2UHWx-AGRk", "1wxppiBVNo2r2dRGfdN7oJYicWePbbHrhU-spQCf9U7I",
		"1q6QaHEziPecGDd7tDgvvCR4bg_J6VENI9ehGrUylx1U", "1dqHHiR40Ref4QkqhJxU1n1ttXeHzyXgEl6VSVIoUbbA", "1842Go0HUVGEV44J_qvnzbN1qiX707NrIqA0gjJx6yrg", "1oeBw7XAMlz_gR57E3kDvJYV3-xy6akgULdntKj5IyOg",
		"1nFl47XN3seedUOzFlo3hUCrDFHLNpcVph0H3nk-yaak", "18Y19f4d3WHYYAq1IMkhUZbS-C4v2MmKkJv5pabXhAWM", "1TggB6guUTWfe6PhVoxHmLiSgaybpEZM0XFUuo5psBkE", "1QQ7_7PUWaT5cR1Af0CGrNa0c_JPTzNE1y5EG8ciZ8YE",
		"17iJ0SPpGSVQKr-RigpJZa9EE7DWs72__c7tha2PsytA", "1WdAM3ZGOJJszy1C0XFEtyG7i3yOGuL1fC97p-cTvZAU", "1MD-UcWfw8Z-lhXBSkrBPkl7-6Q4GBW1KEL4gUDd6pOw", "16vwyTal7ikUlxq5HDlKvyOz1CStaNjh-DKGp79Al6NE",
		"16z3zl-ik8WYPbMgQrBAn7tz1tgtK48adKWxem4PmsoE", "1CDiLmBMh64_nUczpIGb0DRvEwN7dK0v-AQ1KXTm6J2w", "1sCGWsrBVgFbUOjm5JA4u425d5-yPuU2r1aaDGgr-f24", "18Meaq87Ck1Zm_4wOyVrBRWIZUdKaYe4Taj7Gcddpj34",
		"1ewEk5Vi0wFHTIHfsgL2rjljP5Qe2LK--V9x4vwOgtwI", "1vp8Z74I9y_OH0DKLu6FEY2_hs3TvFOqSvKJaKAiv87E", "1H136nbYJZge9GcCJ5m86g4JcGLou-SCfPTQvGeT_McA", "1DusAkoStlVZ-Zu4wGFAuzFJNIHN47ohOK4hVDwxMzbw",
		"1CIDQgAOWT4_aIcCwq4P1AOv7W9koAcupQ90H1o76pVo", "1ygHt7sghAG4GNJuvKuUEcueO1BrtfiUFN-QuyO79ad4", "1dgAaZ3HSq5-Kj4mdPEW-ovas7P4VOUYkXQIIacoTHoo", "1Y_lw9YmehaWLn-IDs9fXgpdLgbEGeHL3FOyeTxl86ew",
		"1KNwqVDR7qQ6sqpP_vJFN110th7QKsWPRA8B3JpZZ99Q", "18mxTBDTFxSdmBM7xir1S_hhPnFezRWB0zLYa9eh_OVg", "1bt7YLhw49THgy748zVLdigO-UKeOr5CyC91ZSvWvAVM", "1yS_4R86jW1cwaHUp23UkMLIjGrMb6MN42vNBbHl-3CE",
		"1NiF677nPlXmWl3u5ZJuxJjkArQbavvTrvZiOl6u2eHg", "1rh5yp4EB1Nfk-YOQpvjPazWn4Wnubn2eUa0vn_K91q4", "1ArP38u-L_VnaUSBdRxleGt7rlU8PSWLJhzJ9hRCkf8U", "1ujXrHDjrGFnD5-uFv3KD4sIgy_ZAlN5TkbkjdoaO3VY",
		"1yo0EejD1g6Np92ad8B2ZwRPjGiW0lSbIAwqW-ZPXwh0", "109sWUkfqZHt0NAW8tt4FQCIflEN-usIBgA4EZiIIUqo", "1_QyLKrEiLid9QupWle-Cwx_AhF8Uu5kz_xzYFosVPZg", "1ZGoTqYnL1n9C_8sRzhXVoccU6iUYPS-S3AWNuFPbJyg",
		"1Si1V6sqqPcDAsRIsBd-seH4AQHnI8lwPLFx4OXFi7K8", "1xy0b090uVuDizYgpXe3Z-LUIdTZ8AlVe9KkJmXoB3tE", "1R4E9xsIHiORRdUBD7qqvnWw2Y7NGRETaN4gUGryz4aE", "1jAVi-klxRf7tTCchV4271MzghdJltP6bb_EVL1fyN-I",
		"1_77bLWWtHUHKkiwYoIqJGoOW7rTpGpCUCMdhj1eCERI", "1gUQZue8ISmwGBbawh4auQ2zb6-O2smWRT-bSGy4oxrA", "1rWU5QFMW_cVmsSvDSRuWh3piWk7hINn-dR11Bv_0-pk", "1_BYswiuoWd2YLiJ3bYLCbIpxWlbYkQfhfZev4_ieU-M",
		"1ZsgpGvX_1iFAzzvCjX5dEBiZA_l1dZW-08j65DYTm0c", "1YPPJ3Qv4frkBYLZzI3oocF5R-YEfuly-bbQcBOWKXrE", "1qBfZ3_eU1bdIjxqyWgct-hvQvkLwB9g6EarLMT-RhCA", "1295iuFaxht8HfiaglgS7UscKAwkgpjwpjldr8tkKBbc",
		"1vphJ_OBLxMdKOkYtm78QnUwAMRCnsZwncb62guOFIe4", "19-d2mEOT-AvfJk4NTUyBngEd3uoaphsAOiiq3L3RYiA", "1opQb-_fLy3MpcRGCt82PY7gxT9-PtxFB5nds9Fh0yXs", "1FwoGFMlOwm4igs-soznJ3lvP9IxVWVZsVf4VrL7i88A",
		"1m6a_eT9h34RE_6aAN6_cWezqAwKhFqk0CkElJ-iX-eQ", "18ZmhgGWoUkrFsBXgG4NYgX-Z9D0_7Fa3R-8ADmgfZ_I", "1jlckU47lRAB471qmQTGNe4JMqOb83Jpw73enjz-SfO8", "1nVRK_uqKpVHYXjkozQRLuK8h3Ipdbh1SzXoIt92Vzic",
		"1opCfdDDlYTxZX0mVh-0UyKeKH7vF3jLeW98WdUA6t1Y", "1TBnkuHM3_qktELXVnzVptNii6QJvN5oCN9H4LbAHRkc", "1OUV2p1jWN9mkZhsVAVVjUR5tIAoffJajIV8MEajX8hs", "1xonCsWjavJATyjX5lWCvoczl4wxRq4apvo9i1sS8k3U",
		"1LYulgQlh8eKBvcXfplbL9grYyYtpWm4ql7dOldvJ_Zc", "1ppdeJBSHDjOnnGsxl1jZLyOBTj8StqqINMhOAgEzKh8", "1RQDGWYYGwCmIa8RNfCd9LmIkPN0aarDZTzy9cAsM6L8", "1OBjbY1lSe0QeXA4mWi8eBITCBjWx0DLh1lX3NvwMwDM",
		"1VYLhBd1xh77kdhxgilxGRE2g5XpszIVM1VP-Sab-7Yw", "1ruS7wHSBeSs6wSnIzjVoHl7ISzqaw3txq7ODvpzz38k", "1Zju_PpNRNy3qjKs-_brr28OfmFDcicynjCKHAPoVBxE", "1aaHEpGSP8ZtRlNiStwttDnyUv12FOh6y86KfLIOjIQY",
		"1PYApk2Yxw63i6y3R210XbZ9fXLKKkX1g5d0MxWyTW3o", "17GOo7ZIE7UgGcRlMYXQu1zZcPskyYqUTkpWctA_b6ZU", "1rb7Ys5InInp-TYcFIwmWZ-_C8kqOE5lALgJZst1p8UM", "1hzJHbloNkbKRHySu56k9fSEmpFFeIhHJlUdHRgCROoI",
		"1irHuX1ZKsYpsjjaDxuUfi_o4JVL4TAM35rufG-0UjGQ", "10yrR9wzDAUdRzUsy7wJ8hS4uJV3z-hGla5NRNBgl6oE", "1oSOSPKBOcIJv_bA8IY9wzx-Cp8FCRcJH-ckNNg8SKw8", "1yaGojYoxNFShLOTcrUzcduVlRcIBD1hh-V8LpMElUWg",
		"1xf96tk8XQ5gco6k8evxxrdbg_b9_EbmFtAVYc-LuhVg", "1KRj3v-2wBeO9IJYaqufF9TjzWQ1-Si8fGaWlnQP9BiQ", "1vIq5_YBiZF9p4TNCa7fi_3JGJ4116CrnWmZNB5uqEbw", "1TgXtWUNIcuQa_QYcThJpW-0vAxB82eOzbz1qGAJKymQ",
		"1TLIwBboooSG8qECGACqUnSa3J63VciE9ndCC-o2tbZ8", "1ACH_viO5Jw-as_5857ssmHp8VeJfTlb5O_8M-7Ai0A0", "1G19LO6V429E2X7O_VFYlqOSTOvZY1ilUq8atECIJHcw", "13Xs5PBwK7u_Zi7hrlfqEon4szBM6HTUMeNXEfZBGsSE",
		"17emaF0JmQxx-l6X4sEDvk2dL-lhTpZaoHvcGirhcizs", "1Khycip6LvADcClgGmsdpXi42wGUsLvVVL5Uf6VgTQdE", "1Xqfv2MNOdCh4QiaJQlhRnSGW6FKSgPP09qkRyzMFxZo", "1pQRb9QupZtdQSPYpSGJbKYZBOgnPe83nL5MWT4HqHxc",
		"1R7v3i3TYWKNHFPRiA3bryLalVn50QCzOi_f8mI67NA8", "1C5iYcDMHCfajIUevaAIAD6ILDUzd7ps9Ezju3j9fpsc", "1NPC3JesGwXi4vizmQIWdGulZDBBPUEnbPlzvxzWjWlU", "1fylgaRE1yF_DpXKq-t_2Z6KFAZ2apwyjc1czckmk8-g",
		"1FP_XEA-RPuoFx_wCtr3-XRLKfsN1ZRWPIGE9iY8TjEU", "1jsSmsI-EC9-FAYr89tTi6MJPVWqrQLyVdN6-cX6_SU4", "1Fa431lAssVU4aOy5yu3T6B2JjynWSLUugD1qE-sPoCE", "1ynG1XRmlNvZLRhMo48uJafHIBLDcCa_i-VeXWLrvQvY",
	}

	spreadsheetIDs = spreadsheetIDs[165:]

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			log.Printf("Spreadsheets.Get: %v", err)
			continue
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)

		for _, namedRange := range spreadsheet.NamedRanges {
			namedRange := namedRange
			if strings.Contains(namedRange.Name, "_blocked_") {
				batch.WithRequest(
					&sheets.Request{
						AddProtectedRange: &sheets.AddProtectedRangeRequest{
							ProtectedRange: &sheets.ProtectedRange{
								Description: namedRange.Name,
								Editors: &sheets.Editors{
									Users: []string{
										"qaztrade.export@gmail.com",
									},
								},
								Range: &sheets.GridRange{
									SheetId:          namedRange.Range.SheetId,
									StartColumnIndex: namedRange.Range.StartColumnIndex,
									EndColumnIndex:   namedRange.Range.EndColumnIndex,
								},
							},
						},
					},
				)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return fmt.Errorf("batch.Do: %w", err)
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}
