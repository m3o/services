module github.com/kytra-app/helpers/textenhancer/cli

replace github.com/kytra-app/helpers/textenhancer => ../

replace github.com/kytra-app/helpers/microgorm => ../../microgorm

replace github.com/kytra-app/helpers/passwordhasher => ../../passwordhasher

replace github.com/kytra-app/users-srv => ../../../users-srv

replace github.com/kytra-app/stocks-srv => ../../../stocks-srv

go 1.12

require github.com/kytra-app/helpers/textenhancer v0.0.0-00010101000000-000000000000
