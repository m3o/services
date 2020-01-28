package handler

import (
	"context"

	proto "github.com/kytra-app/welcome-api/proto"
)

// GetConfig returns the configuration for the "Welcome" interfaces
func (h Handler) GetConfig(ctx context.Context, req *proto.GetConfigRequest, rsp *proto.GetConfigResponse) error {
	rsp.DefaultIndustryAllocations = industryAllocations
	rsp.CountryCodes = countryCodes
	return nil
}

var industryAllocations = []*proto.IndustryAllocation{
	{
		Industry:          "Information Technology",
		DefaultPercentage: 21.36,
	},
	{
		Industry:          "Financials",
		DefaultPercentage: 14.26,
	},
	{
		Industry:          "Health Care",
		DefaultPercentage: 12.85,
	},
	{
		Industry:          "Consumer Discretionary",
		DefaultPercentage: 10.07,
	},
	{
		Industry:          "Communication Services",
		DefaultPercentage: 9.94,
	},
	{
		Industry:          "Industrials",
		DefaultPercentage: 9.29,
	},
	{
		Industry:          "Consumer Staples",
		DefaultPercentage: 7.21,
	},
	{
		Industry:          "Energy",
		DefaultPercentage: .52,
	},
	{
		Industry:          "Utilities",
		DefaultPercentage: 3.48,
	},
	{
		Industry:          "Real Estate",
		DefaultPercentage: 3.28,
	},
	{
		Industry:          "Materials",
		DefaultPercentage: 3.08,
	},
}

var countryCodes = []*proto.CountryCode{
	{
		Country:  "Israel",
		DialCode: "+972",
		Code:     "IL",
	},
	{
		Country:  "Afghanistan",
		DialCode: "+93",
		Code:     "AF",
	},
	{
		Country:  "Albania",
		DialCode: "+355",
		Code:     "AL",
	},
	{
		Country:  "Algeria",
		DialCode: "+213",
		Code:     "DZ",
	},
	{
		Country:  "AmericanSamoa",
		DialCode: "+1 684",
		Code:     "AS",
	},
	{
		Country:  "Andorra",
		DialCode: "+376",
		Code:     "AD",
	},
	{
		Country:  "Angola",
		DialCode: "+244",
		Code:     "AO",
	},
	{
		Country:  "Anguilla",
		DialCode: "+1 264",
		Code:     "AI",
	},
	{
		Country:  "Antigua and Barbuda",
		DialCode: "+1268",
		Code:     "AG",
	},
	{
		Country:  "Argentina",
		DialCode: "+54",
		Code:     "AR",
	},
	{
		Country:  "Armenia",
		DialCode: "+374",
		Code:     "AM",
	},
	{
		Country:  "Aruba",
		DialCode: "+297",
		Code:     "AW",
	},
	{
		Country:  "Australia",
		DialCode: "+61",
		Code:     "AU",
	},
	{
		Country:  "Austria",
		DialCode: "+43",
		Code:     "AT",
	},
	{
		Country:  "Azerbaijan",
		DialCode: "+994",
		Code:     "AZ",
	},
	{
		Country:  "Bahamas",
		DialCode: "+1 242",
		Code:     "BS",
	},
	{
		Country:  "Bahrain",
		DialCode: "+973",
		Code:     "BH",
	},
	{
		Country:  "Bangladesh",
		DialCode: "+880",
		Code:     "BD",
	},
	{
		Country:  "Barbados",
		DialCode: "+1 246",
		Code:     "BB",
	},
	{
		Country:  "Belarus",
		DialCode: "+375",
		Code:     "BY",
	},
	{
		Country:  "Belgium",
		DialCode: "+32",
		Code:     "BE",
	},
	{
		Country:  "Belize",
		DialCode: "+501",
		Code:     "BZ",
	},
	{
		Country:  "Benin",
		DialCode: "+229",
		Code:     "BJ",
	},
	{
		Country:  "Bermuda",
		DialCode: "+1 441",
		Code:     "BM",
	},
	{
		Country:  "Bhutan",
		DialCode: "+975",
		Code:     "BT",
	},
	{
		Country:  "Bosnia and Herzegovina",
		DialCode: "+387",
		Code:     "BA",
	},
	{
		Country:  "Botswana",
		DialCode: "+267",
		Code:     "BW",
	},
	{
		Country:  "Brazil",
		DialCode: "+55",
		Code:     "BR",
	},
	{
		Country:  "British Indian Ocean Territory",
		DialCode: "+246",
		Code:     "IO",
	},
	{
		Country:  "Bulgaria",
		DialCode: "+359",
		Code:     "BG",
	},
	{
		Country:  "Burkina Faso",
		DialCode: "+226",
		Code:     "BF",
	},
	{
		Country:  "Burundi",
		DialCode: "+257",
		Code:     "BI",
	},
	{
		Country:  "Cambodia",
		DialCode: "+855",
		Code:     "KH",
	},
	{
		Country:  "Cameroon",
		DialCode: "+237",
		Code:     "CM",
	},
	{
		Country:  "Canada",
		DialCode: "+1",
		Code:     "CA",
	},
	{
		Country:  "Cape Verde",
		DialCode: "+238",
		Code:     "CV",
	},
	{
		Country:  "Cayman Islands",
		DialCode: "+ 345",
		Code:     "KY",
	},
	{
		Country:  "Central African Republic",
		DialCode: "+236",
		Code:     "CF",
	},
	{
		Country:  "Chad",
		DialCode: "+235",
		Code:     "TD",
	},
	{
		Country:  "Chile",
		DialCode: "+56",
		Code:     "CL",
	},
	{
		Country:  "China",
		DialCode: "+86",
		Code:     "CN",
	},
	{
		Country:  "Christmas Island",
		DialCode: "+61",
		Code:     "CX",
	},
	{
		Country:  "Colombia",
		DialCode: "+57",
		Code:     "CO",
	},
	{
		Country:  "Comoros",
		DialCode: "+269",
		Code:     "KM",
	},
	{
		Country:  "Congo",
		DialCode: "+242",
		Code:     "CG",
	},
	{
		Country:  "Cook Islands",
		DialCode: "+682",
		Code:     "CK",
	},
	{
		Country:  "Costa Rica",
		DialCode: "+506",
		Code:     "CR",
	},
	{
		Country:  "Croatia",
		DialCode: "+385",
		Code:     "HR",
	},
	{
		Country:  "Cuba",
		DialCode: "+53",
		Code:     "CU",
	},
	{
		Country:  "Cyprus",
		DialCode: "+537",
		Code:     "CY",
	},
	{
		Country:  "Czech Republic",
		DialCode: "+420",
		Code:     "CZ",
	},
	{
		Country:  "Denmark",
		DialCode: "+45",
		Code:     "DK",
	},
	{
		Country:  "Djibouti",
		DialCode: "+253",
		Code:     "DJ",
	},
	{
		Country:  "Dominica",
		DialCode: "+1 767",
		Code:     "DM",
	},
	{
		Country:  "Dominican Republic",
		DialCode: "+1 849",
		Code:     "DO",
	},
	{
		Country:  "Ecuador",
		DialCode: "+593",
		Code:     "EC",
	},
	{
		Country:  "Egypt",
		DialCode: "+20",
		Code:     "EG",
	},
	{
		Country:  "El Salvador",
		DialCode: "+503",
		Code:     "SV",
	},
	{
		Country:  "Equatorial Guinea",
		DialCode: "+240",
		Code:     "GQ",
	},
	{
		Country:  "Eritrea",
		DialCode: "+291",
		Code:     "ER",
	},
	{
		Country:  "Estonia",
		DialCode: "+372",
		Code:     "EE",
	},
	{
		Country:  "Ethiopia",
		DialCode: "+251",
		Code:     "ET",
	},
	{
		Country:  "Faroe Islands",
		DialCode: "+298",
		Code:     "FO",
	},
	{
		Country:  "Fiji",
		DialCode: "+679",
		Code:     "FJ",
	},
	{
		Country:  "Finland",
		DialCode: "+358",
		Code:     "FI",
	},
	{
		Country:  "France",
		DialCode: "+33",
		Code:     "FR",
	},
	{
		Country:  "French Guiana",
		DialCode: "+594",
		Code:     "GF",
	},
	{
		Country:  "French Polynesia",
		DialCode: "+689",
		Code:     "PF",
	},
	{
		Country:  "Gabon",
		DialCode: "+241",
		Code:     "GA",
	},
	{
		Country:  "Gambia",
		DialCode: "+220",
		Code:     "GM",
	},
	{
		Country:  "Georgia",
		DialCode: "+995",
		Code:     "GE",
	},
	{
		Country:  "Germany",
		DialCode: "+49",
		Code:     "DE",
	},
	{
		Country:  "Ghana",
		DialCode: "+233",
		Code:     "GH",
	},
	{
		Country:  "Gibraltar",
		DialCode: "+350",
		Code:     "GI",
	},
	{
		Country:  "Greece",
		DialCode: "+30",
		Code:     "GR",
	},
	{
		Country:  "Greenland",
		DialCode: "+299",
		Code:     "GL",
	},
	{
		Country:  "Grenada",
		DialCode: "+1 473",
		Code:     "GD",
	},
	{
		Country:  "Guadeloupe",
		DialCode: "+590",
		Code:     "GP",
	},
	{
		Country:  "Guam",
		DialCode: "+1 671",
		Code:     "GU",
	},
	{
		Country:  "Guatemala",
		DialCode: "+502",
		Code:     "GT",
	},
	{
		Country:  "Guinea",
		DialCode: "+224",
		Code:     "GN",
	},
	{
		Country:  "Guinea-Bissau",
		DialCode: "+245",
		Code:     "GW",
	},
	{
		Country:  "Guyana",
		DialCode: "+595",
		Code:     "GY",
	},
	{
		Country:  "Haiti",
		DialCode: "+509",
		Code:     "HT",
	},
	{
		Country:  "Honduras",
		DialCode: "+504",
		Code:     "HN",
	},
	{
		Country:  "Hungary",
		DialCode: "+36",
		Code:     "HU",
	},
	{
		Country:  "Iceland",
		DialCode: "+354",
		Code:     "IS",
	},
	{
		Country:  "India",
		DialCode: "+91",
		Code:     "IN",
	},
	{
		Country:  "Indonesia",
		DialCode: "+62",
		Code:     "ID",
	},
	{
		Country:  "Iraq",
		DialCode: "+964",
		Code:     "IQ",
	},
	{
		Country:  "Ireland",
		DialCode: "+353",
		Code:     "IE",
	},
	{
		Country:  "Israel",
		DialCode: "+972",
		Code:     "IL",
	},
	{
		Country:  "Italy",
		DialCode: "+39",
		Code:     "IT",
	},
	{
		Country:  "Jamaica",
		DialCode: "+1 876",
		Code:     "JM",
	},
	{
		Country:  "Japan",
		DialCode: "+81",
		Code:     "JP",
	},
	{
		Country:  "Jordan",
		DialCode: "+962",
		Code:     "JO",
	},
	{
		Country:  "Kazakhstan",
		DialCode: "+7 7",
		Code:     "KZ",
	},
	{
		Country:  "Kenya",
		DialCode: "+254",
		Code:     "KE",
	},
	{
		Country:  "Kiribati",
		DialCode: "+686",
		Code:     "KI",
	},
	{
		Country:  "Kuwait",
		DialCode: "+965",
		Code:     "KW",
	},
	{
		Country:  "Kyrgyzstan",
		DialCode: "+996",
		Code:     "KG",
	},
	{
		Country:  "Latvia",
		DialCode: "+371",
		Code:     "LV",
	},
	{
		Country:  "Lebanon",
		DialCode: "+961",
		Code:     "LB",
	},
	{
		Country:  "Lesotho",
		DialCode: "+266",
		Code:     "LS",
	},
	{
		Country:  "Liberia",
		DialCode: "+231",
		Code:     "LR",
	},
	{
		Country:  "Liechtenstein",
		DialCode: "+423",
		Code:     "LI",
	},
	{
		Country:  "Lithuania",
		DialCode: "+370",
		Code:     "LT",
	},
	{
		Country:  "Luxembourg",
		DialCode: "+352",
		Code:     "LU",
	},
	{
		Country:  "Madagascar",
		DialCode: "+261",
		Code:     "MG",
	},
	{
		Country:  "Malawi",
		DialCode: "+265",
		Code:     "MW",
	},
	{
		Country:  "Malaysia",
		DialCode: "+60",
		Code:     "MY",
	},
	{
		Country:  "Maldives",
		DialCode: "+960",
		Code:     "MV",
	},
	{
		Country:  "Mali",
		DialCode: "+223",
		Code:     "ML",
	},
	{
		Country:  "Malta",
		DialCode: "+356",
		Code:     "MT",
	},
	{
		Country:  "Marshall Islands",
		DialCode: "+692",
		Code:     "MH",
	},
	{
		Country:  "Martinique",
		DialCode: "+596",
		Code:     "MQ",
	},
	{
		Country:  "Mauritania",
		DialCode: "+222",
		Code:     "MR",
	},
	{
		Country:  "Mauritius",
		DialCode: "+230",
		Code:     "MU",
	},
	{
		Country:  "Mayotte",
		DialCode: "+262",
		Code:     "YT",
	},
	{
		Country:  "Mexico",
		DialCode: "+52",
		Code:     "MX",
	},
	{
		Country:  "Monaco",
		DialCode: "+377",
		Code:     "MC",
	},
	{
		Country:  "Mongolia",
		DialCode: "+976",
		Code:     "MN",
	},
	{
		Country:  "Montenegro",
		DialCode: "+382",
		Code:     "ME",
	},
	{
		Country:  "Montserrat",
		DialCode: "+1664",
		Code:     "MS",
	},
	{
		Country:  "Morocco",
		DialCode: "+212",
		Code:     "MA",
	},
	{
		Country:  "Myanmar",
		DialCode: "+95",
		Code:     "MM",
	},
	{
		Country:  "Namibia",
		DialCode: "+264",
		Code:     "NA",
	},
	{
		Country:  "Nauru",
		DialCode: "+674",
		Code:     "NR",
	},
	{
		Country:  "Nepal",
		DialCode: "+977",
		Code:     "NP",
	},
	{
		Country:  "Netherlands",
		DialCode: "+31",
		Code:     "NL",
	},
	{
		Country:  "Netherlands Antilles",
		DialCode: "+599",
		Code:     "AN",
	},
	{
		Country:  "New Caledonia",
		DialCode: "+687",
		Code:     "NC",
	},
	{
		Country:  "New Zealand",
		DialCode: "+64",
		Code:     "NZ",
	},
	{
		Country:  "Nicaragua",
		DialCode: "+505",
		Code:     "NI",
	},
	{
		Country:  "Niger",
		DialCode: "+227",
		Code:     "NE",
	},
	{
		Country:  "Nigeria",
		DialCode: "+234",
		Code:     "NG",
	},
	{
		Country:  "Niue",
		DialCode: "+683",
		Code:     "NU",
	},
	{
		Country:  "Norfolk Island",
		DialCode: "+672",
		Code:     "NF",
	},
	{
		Country:  "Northern Mariana Islands",
		DialCode: "+1 670",
		Code:     "MP",
	},
	{
		Country:  "Norway",
		DialCode: "+47",
		Code:     "NO",
	},
	{
		Country:  "Oman",
		DialCode: "+968",
		Code:     "OM",
	},
	{
		Country:  "Pakistan",
		DialCode: "+92",
		Code:     "PK",
	},
	{
		Country:  "Palau",
		DialCode: "+680",
		Code:     "PW",
	},
	{
		Country:  "Panama",
		DialCode: "+507",
		Code:     "PA",
	},
	{
		Country:  "Papua New Guinea",
		DialCode: "+675",
		Code:     "PG",
	},
	{
		Country:  "Paraguay",
		DialCode: "+595",
		Code:     "PY",
	},
	{
		Country:  "Peru",
		DialCode: "+51",
		Code:     "PE",
	},
	{
		Country:  "Philippines",
		DialCode: "+63",
		Code:     "PH",
	},
	{
		Country:  "Poland",
		DialCode: "+48",
		Code:     "PL",
	},
	{
		Country:  "Portugal",
		DialCode: "+351",
		Code:     "PT",
	},
	{
		Country:  "Puerto Rico",
		DialCode: "+1 939",
		Code:     "PR",
	},
	{
		Country:  "Qatar",
		DialCode: "+974",
		Code:     "QA",
	},
	{
		Country:  "Romania",
		DialCode: "+40",
		Code:     "RO",
	},
	{
		Country:  "Rwanda",
		DialCode: "+250",
		Code:     "RW",
	},
	{
		Country:  "Samoa",
		DialCode: "+685",
		Code:     "WS",
	},
	{
		Country:  "San Marino",
		DialCode: "+378",
		Code:     "SM",
	},
	{
		Country:  "Saudi Arabia",
		DialCode: "+966",
		Code:     "SA",
	},
	{
		Country:  "Senegal",
		DialCode: "+221",
		Code:     "SN",
	},
	{
		Country:  "Serbia",
		DialCode: "+381",
		Code:     "RS",
	},
	{
		Country:  "Seychelles",
		DialCode: "+248",
		Code:     "SC",
	},
	{
		Country:  "Sierra Leone",
		DialCode: "+232",
		Code:     "SL",
	},
	{
		Country:  "Singapore",
		DialCode: "+65",
		Code:     "SG",
	},
	{
		Country:  "Slovakia",
		DialCode: "+421",
		Code:     "SK",
	},
	{
		Country:  "Slovenia",
		DialCode: "+386",
		Code:     "SI",
	},
	{
		Country:  "Solomon Islands",
		DialCode: "+677",
		Code:     "SB",
	},
	{
		Country:  "South Africa",
		DialCode: "+27",
		Code:     "ZA",
	},
	{
		Country:  "South Georgia and the South Sandwich Islands",
		DialCode: "+500",
		Code:     "GS",
	},
	{
		Country:  "Spain",
		DialCode: "+34",
		Code:     "ES",
	},
	{
		Country:  "Sri Lanka",
		DialCode: "+94",
		Code:     "LK",
	},
	{
		Country:  "Sudan",
		DialCode: "+249",
		Code:     "SD",
	},
	{
		Country:  "Suriname",
		DialCode: "+597",
		Code:     "SR",
	},
	{
		Country:  "Swaziland",
		DialCode: "+268",
		Code:     "SZ",
	},
	{
		Country:  "Sweden",
		DialCode: "+46",
		Code:     "SE",
	},
	{
		Country:  "Switzerland",
		DialCode: "+41",
		Code:     "CH",
	},
	{
		Country:  "Tajikistan",
		DialCode: "+992",
		Code:     "TJ",
	},
	{
		Country:  "Thailand",
		DialCode: "+66",
		Code:     "TH",
	},
	{
		Country:  "Togo",
		DialCode: "+228",
		Code:     "TG",
	},
	{
		Country:  "Tokelau",
		DialCode: "+690",
		Code:     "TK",
	},
	{
		Country:  "Tonga",
		DialCode: "+676",
		Code:     "TO",
	},
	{
		Country:  "Trinidad and Tobago",
		DialCode: "+1 868",
		Code:     "TT",
	},
	{
		Country:  "Tunisia",
		DialCode: "+216",
		Code:     "TN",
	},
	{
		Country:  "Turkey",
		DialCode: "+90",
		Code:     "TR",
	},
	{
		Country:  "Turkmenistan",
		DialCode: "+993",
		Code:     "TM",
	},
	{
		Country:  "Turks and Caicos Islands",
		DialCode: "+1 649",
		Code:     "TC",
	},
	{
		Country:  "Tuvalu",
		DialCode: "+688",
		Code:     "TV",
	},
	{
		Country:  "Uganda",
		DialCode: "+256",
		Code:     "UG",
	},
	{
		Country:  "Ukraine",
		DialCode: "+380",
		Code:     "UA",
	},
	{
		Country:  "United Arab Emirates",
		DialCode: "+971",
		Code:     "AE",
	},
	{
		Country:  "United Kingdom",
		DialCode: "+44",
		Code:     "GB",
	},
	{
		Country:  "United States",
		DialCode: "+1",
		Code:     "US",
	},
	{
		Country:  "Uruguay",
		DialCode: "+598",
		Code:     "UY",
	},
	{
		Country:  "Uzbekistan",
		DialCode: "+998",
		Code:     "UZ",
	},
	{
		Country:  "Vanuatu",
		DialCode: "+678",
		Code:     "VU",
	},
	{
		Country:  "Wallis and Futuna",
		DialCode: "+681",
		Code:     "WF",
	},
	{
		Country:  "Yemen",
		DialCode: "+967",
		Code:     "YE",
	},
	{
		Country:  "Zambia",
		DialCode: "+260",
		Code:     "ZM",
	},
	{
		Country:  "Zimbabwe",
		DialCode: "+263",
		Code:     "ZW",
	},
	{
		Country:  "land Islands",
		DialCode: "",
		Code:     "AX",
	},
	{
		Country:  "Bolivia, Plurinational State of",
		DialCode: "+591",
		Code:     "BO",
	},
	{
		Country:  "Brunei Darussalam",
		DialCode: "+673",
		Code:     "BN",
	},
	{
		Country:  "Cocos (Keeling) Islands",
		DialCode: "+61",
		Code:     "CC",
	},
	{
		Country:  "Congo, The Democratic Republic of the",
		DialCode: "+243",
		Code:     "CD",
	},
	{
		Country:  "Cote d'Ivoire",
		DialCode: "+225",
		Code:     "CI",
	},
	{
		Country:  "Falkland Islands (Malvinas)",
		DialCode: "+500",
		Code:     "FK",
	},
	{
		Country:  "Guernsey",
		DialCode: "+44",
		Code:     "GG",
	},
	{
		Country:  "Holy See (Vatican City State)",
		DialCode: "+379",
		Code:     "VA",
	},
	{
		Country:  "Hong Kong",
		DialCode: "+852",
		Code:     "HK",
	},
	{
		Country:  "Iran, Islamic Republic of",
		DialCode: "+98",
		Code:     "IR",
	},
	{
		Country:  "Isle of Man",
		DialCode: "+44",
		Code:     "IM",
	},
	{
		Country:  "Jersey",
		DialCode: "+44",
		Code:     "JE",
	},
	{
		Country:  "Korea, Democratic People's Republic of",
		DialCode: "+850",
		Code:     "KP",
	},
	{
		Country:  "Korea, Republic of",
		DialCode: "+82",
		Code:     "KR",
	},
	{
		Country:  "Lao People's Democratic Republic",
		DialCode: "+856",
		Code:     "LA",
	},
	{
		Country:  "Libyan Arab Jamahiriya",
		DialCode: "+218",
		Code:     "LY",
	},
	{
		Country:  "Macao",
		DialCode: "+853",
		Code:     "MO",
	},
	{
		Country:  "Macedonia, The Former Yugoslav Republic of",
		DialCode: "+389",
		Code:     "MK",
	},
	{
		Country:  "Micronesia, Federated States of",
		DialCode: "+691",
		Code:     "FM",
	},
	{
		Country:  "Moldova, Republic of",
		DialCode: "+373",
		Code:     "MD",
	},
	{
		Country:  "Mozambique",
		DialCode: "+258",
		Code:     "MZ",
	},
	{
		Country:  "Palestinian Territory, Occupied",
		DialCode: "+970",
		Code:     "PS",
	},
	{
		Country:  "Pitcairn",
		DialCode: "+872",
		Code:     "PN",
	},
	{
		Country:  "Réunion",
		DialCode: "+262",
		Code:     "RE",
	},
	{
		Country:  "Russia",
		DialCode: "+7",
		Code:     "RU",
	},
	{
		Country:  "Saint Barthélemy",
		DialCode: "+590",
		Code:     "BL",
	},
	{
		Country:  "Saint Helena, Ascension and Tristan Da Cunha",
		DialCode: "+290",
		Code:     "SH",
	},
	{
		Country:  "Saint Kitts and Nevis",
		DialCode: "+1 869",
		Code:     "KN",
	},
	{
		Country:  "Saint Lucia",
		DialCode: "+1 758",
		Code:     "LC",
	},
	{
		Country:  "Saint Martin",
		DialCode: "+590",
		Code:     "MF",
	},
	{
		Country:  "Saint Pierre and Miquelon",
		DialCode: "+508",
		Code:     "PM",
	},
	{
		Country:  "Saint Vincent and the Grenadines",
		DialCode: "+1 784",
		Code:     "VC",
	},
	{
		Country:  "Sao Tome and Principe",
		DialCode: "+239",
		Code:     "ST",
	},
	{
		Country:  "Somalia",
		DialCode: "+252",
		Code:     "SO",
	},
	{
		Country:  "Svalbard and Jan Mayen",
		DialCode: "+47",
		Code:     "SJ",
	},
	{
		Country:  "Syrian Arab Republic",
		DialCode: "+963",
		Code:     "SY",
	},
	{
		Country:  "Taiwan, Province of China",
		DialCode: "+886",
		Code:     "TW",
	},
	{
		Country:  "Tanzania, United Republic of",
		DialCode: "+255",
		Code:     "TZ",
	},
	{
		Country:  "Timor-Leste",
		DialCode: "+670",
		Code:     "TL",
	},
	{
		Country:  "Venezuela, Bolivarian Republic of",
		DialCode: "+58",
		Code:     "VE",
	},
	{
		Country:  "Viet Nam",
		DialCode: "+84",
		Code:     "VN",
	},
	{
		Country:  "Virgin Islands, British",
		DialCode: "+1 284",
		Code:     "VG",
	},
	{
		Country:  "Virgin Islands, U.S.",
		DialCode: "+1 340",
		Code:     "VI",
	},
}
