-include gomk/main.mk
-include local/Makefile

clean: clean-default
ifeq ($(unameS),windows)
ifneq ($(wildcard resource_windows*.syso),)
	@remove-item -force ./cmd/xgo/resource_windows*.syso
endif
else
	@rm -f ./cmd/xgo/resource_windows*.syso
endif
