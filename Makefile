BINDIR = /usr/local/bin

help:          ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:       ## Install Target
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o ${BINDIR}/_kxd_prompt
	cp scripts/_kxd ${BINDIR}/_kxd
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
	@echo "        __ __   _  __    ____   "
	@echo "       / //_/  | |/ /   / __ \  "
	@echo "      / ,<     |   /   / / / /  "
	@echo "     / /| |   /   |   / /_/ /   "
	@echo "    /_/ |_|  /_/|_|  /_____/    "
	@echo "                                "
	@echo "   To Finish Installation add   "
	@echo "                                "
	@echo "   alias kxd=\"source _kxd\"    "
	@echo "                                "
	@echo " to your bash profile or zshrc  "
	@echo "   then open new terminal or    "
	@echo "       source that file         "
	@echo "                                "
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "

uninstall:     ## Uninstall Target
	rm -f ${BINDIR}/_kxd
	rm -f ${BINDIR}/_kxd_prompt
