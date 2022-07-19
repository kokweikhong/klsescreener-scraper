package keys

type BOARD int

const (
	B_MAIN_MARKET BOARD = iota + 1
	B_ACE_MARKET
	B_STRUCTURED_WARRANTS
	B_ETF
	B_BOND_AND_LOAN
	B_LEAP_MARKET
)

type SUB_SECTOR int

const (
	SUB_AGRICULTURAL_PRODUCTS                         SUB_SECTOR = 5
	SUB_AUTO_PARTS                                    SUB_SECTOR = 39
	SUB_AUTOMOTIVE                                    SUB_SECTOR = 7
	SUB_BANKING                                       SUB_SECTOR = 27
	SUB_BOND_FUND                                     SUB_SECTOR = 91
	SUB_BUILDING_MATERIALS                            SUB_SECTOR = 41
	SUB_CHEMICALS                                     SUB_SECTOR = 43
	SUB_CLOSED_END_FUND                               SUB_SECTOR = 1
	SUB_COMMODITY_FUND                                SUB_SECTOR = 93
	SUB_CONSTRUCTION                                  SUB_SECTOR = 3
	SUB_CONSUMER_SERVICES                             SUB_SECTOR = 9
	SUB_CONVENTIONAL_GG                               SUB_SECTOR = 97
	SUB_CONVENTIONAL_MGS                              SUB_SECTOR = 99
	SUB_CONVENTIONAL_PDS                              SUB_SECTOR = 101
	SUB_DIGITAL_SERVICES                              SUB_SECTOR = 67
	SUB_DIVERSIFIED_INDUSTRIALS                       SUB_SECTOR = 45
	SUB_ELECTRICITY                                   SUB_SECTOR = 85
	SUB_ENERGY_INFRASTRUCTURE_EQUIPMENT_AND_SERVICES  SUB_SECTOR = 21
	SUB_EQUITY_FUND                                   SUB_SECTOR = 95
	SUB_FOOD_AND_BEVERAGES                            SUB_SECTOR = 11
	SUB_GAS_WATER_AND_MULTI_UTILITIES                 SUB_SECTOR = 87
	SUB_HEALTH_CARE_EQUIPMENT_AND_SERVICES            SUB_SECTOR = 33
	SUB_HEALTH_CARE_PROVIDERS                         SUB_SECTOR = 35
	SUB_HOUSHOLD_GOODS                                SUB_SECTOR = 13
	SUB_INDUSTRIAL_ENGINEERING                        SUB_SECTOR = 47
	SUB_INDUSTRIAL_MATERIALS_COMPONENTS_AND_EQUIPMENT SUB_SECTOR = 49
	SUB_INDUSTRIAL_SERVICES                           SUB_SECTOR = 51
	SUB_INSURANCE                                     SUB_SECTOR = 29
	SUB_ISLAMIC_GG                                    SUB_SECTOR = 103
	SUB_ISLAMIC_GII                                   SUB_SECTOR = 105
	SUB_ISLAMIC_PDS                                   SUB_SECTOR = 107
	SUB_MEDIA                                         SUB_SECTOR = 75
	SUB_METALS                                        SUB_SECTOR = 53
	SUB_OIL_AND_GAS_PRODUCERS                         SUB_SECTOR = 23
	SUB_OTHER_ENERGY_RESOURCES                        SUB_SECTOR = 25
	SUB_OTHER_FINANCIALS                              SUB_SECTOR = 31
	SUB_PACKAGING_MATERIALS                           SUB_SECTOR = 55
	SUB_PERSONAL_GOODS                                SUB_SECTOR = 15
	SUB_PHARMACEUTICALS                               SUB_SECTOR = 37
	SUB_PLANTATION                                    SUB_SECTOR = 59
	SUB_PROPERTY                                      SUB_SECTOR = 61
	SUB_REAL_ESTATE_INVESTMENT_TRUSTS                 SUB_SECTOR = 63
	SUB_RETAILERS                                     SUB_SECTOR = 17
	SUB_SEMICONDUCTORS                                SUB_SECTOR = 69
	SUB_SOFTWARE                                      SUB_SECTOR = 71
	SUB_SPECIAL_PURPOSE_ACQUISITION_COMPANY           SUB_SECTOR = 65
	SUB_STRUCTURED_WARRANTS                           SUB_SECTOR = 89
	SUB_TECHNOLOGY_EQUIPMENT                          SUB_SECTOR = 73
	SUB_TELECOMMUNICATIONS_EQUIPMENT                  SUB_SECTOR = 77
	SUB_TELECOMMUNICATIONS_SERVICE_PROVIDERS          SUB_SECTOR = 79
	SUB_TRANSPORTATION_AND_LOGISTICS_SERVICES         SUB_SECTOR = 81
	SUB_TRANSPORTATION_EQUIPMENT                      SUB_SECTOR = 83
	SUB_TRAVEL_LEISURE_AND_HOSPITALITY                SUB_SECTOR = 19
	SUB_WOOD_AND_WOOD_PRODUCTS                        SUB_SECTOR = 57
)
