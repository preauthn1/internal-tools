
go-test: ## Run tests
	bash ./scripts/gotest.sh
go-build-init:
	bash ./scripts/gitlab-ci-build-init.sh
go-build-mains:
	bash ./scripts/gitlab-ci-build-mains.sh
go-build-games:
	bash ./scripts/gitlab-ci-build-games.sh
go-download-modules:
	bash ./scripts/gitlab-ci-download-modules.sh

docker-build-games:
	bash ./scripts/gitlab-ci-build-games-docker.sh

build-main:
	bash ./scripts/buildmains.sh
	
build-game:
	bash ./scripts/buildgames.sh

test-game-dfdc:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/dfdc/config/gameserver.dat.example ./internal/game/dfdc/config/gameserver.dat
	cp ./internal/game/dfdc/config/serverinfo.dat.example ./internal/game/dfdc/config/serverinfo.dat
	cp ./internal/game/dfdc/config/slotconfig.dat.example ./internal/game/dfdc/config/slotconfig.dat

	go test ./internal/game/dfdc -coverprofile=.cover/internal.game.dfdc.cover

test-game-tgw:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/candyhouse/config/gameserver.dat.example ./internal/game/candyhouse/config/gameserver.dat
	cp ./internal/game/candyhouse/config/serverinfo.dat.example ./internal/game/candyhouse/config/serverinfo.dat
	cp ./internal/game/candyhouse/config/slotconfig.dat.example ./internal/game/candyhouse/config/slotconfig.dat

	go test ./internal/game/candyhouse -coverprofile=.cover/internal.game.candyhouse.cover

test-game-ald:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/aladdin/config/gameserver.dat.example ./internal/game/aladdin/config/gameserver.dat
	cp ./internal/game/aladdin/config/serverinfo.dat.example ./internal/game/aladdin/config/serverinfo.dat
	cp ./internal/game/aladdin/config/slotconfig.dat.example ./internal/game/aladdin/config/slotconfig.dat

	go test ./internal/game/aladdin -coverprofile=.cover/internal.game.aladdin.cover

test-game-wlzb:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/wlzb/config/gameserver.dat.example ./internal/game/wlzb/config/gameserver.dat
	cp ./internal/game/wlzb/config/serverinfo.dat.example ./internal/game/wlzb/config/serverinfo.dat
	cp ./internal/game/wlzb/config/slotconfig.dat.example ./internal/game/wlzb/config/slotconfig.dat

	go test ./internal/game/wlzb -coverprofile=.cover/internal.game.wlzb.cover

test-game-xmbx:
	mkdir -p .cover/html

	# DB Handler
	go test ./internal/game/xmbx/dbhandler/service -coverprofile=.cover/internal.game.xmbx.dbhandler.service.cover
	go test ./internal/game/xmbx/dbhandler -coverprofile=.cover/internal.game.xmbx.dbhandler.cover

	# Copy config file
	cp ./internal/game/xmbx/config/gameserver.dat.example ./internal/game/xmbx/config/gameserver.dat
	cp ./internal/game/xmbx/config/serverinfo.dat.example ./internal/game/xmbx/config/serverinfo.dat
	cp ./internal/game/xmbx/config/slotconfig.dat.example ./internal/game/xmbx/config/slotconfig.dat

	go test ./internal/game/xmbx -coverprofile=.cover/internal.game.xmbx.cover

test-game-jlbpk:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/jlbpk/config/gameserver.dat.example ./internal/game/jlbpk/config/gameserver.dat
	cp ./internal/game/jlbpk/config/serverinfo.dat.example ./internal/game/jlbpk/config/serverinfo.dat
	cp ./internal/game/jlbpk/config/poker.dat.example ./internal/game/jlbpk/config/poker.dat
	cp ./internal/game/jlbpk/config/syspool.dat.example ./internal/game/jlbpk/config/syspool.dat

	go test ./internal/game/jlbpk -coverprofile=.cover/internal.game.jlbpk.cover

test-game-jjbx:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/jjbx/config/gameserver.dat.example ./internal/game/jjbx/config/gameserver.dat
	cp ./internal/game/jjbx/config/serverinfo.dat.example ./internal/game/jjbx/config/serverinfo.dat
	cp ./internal/game/jjbx/config/slotconfig.dat.example ./internal/game/jjbx/config/slotconfig.dat

	go test ./internal/game/jjbx -coverprofile=.cover/internal.game.jjbx.cover
	
test-game-jdesyd:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/jdesyd/config/gameserver.dat.example ./internal/game/jdesyd/config/gameserver.dat
	cp ./internal/game/jdesyd/config/serverinfo.dat.example ./internal/game/jdesyd/config/serverinfo.dat
	cp ./internal/game/jdesyd/config/poker.dat.example ./internal/game/jdesyd/config/poker.dat
	cp ./internal/game/jdesyd/config/syspool.dat.example ./internal/game/jdesyd/config/syspool.dat

	go test ./internal/game/jdesyd -coverprofile=.cover/internal.game.jdesyd.cover
	
test-game-dczz:
	mkdir -p .cover/html

	# Copy config file
	cp ./internal/game/dczz/config/gameserver.dat.example ./internal/game/dczz/config/gameserver.dat
	cp ./internal/game/dczz/config/serverinfo.dat.example ./internal/game/dczz/config/serverinfo.dat
	cp ./internal/game/dczz/config/dealingdevice.dat.example ./internal/game/dczz/config/dealingdevice.dat

	go test ./internal/game/dczz -coverprofile=.cover/internal.game.dczz.cover
