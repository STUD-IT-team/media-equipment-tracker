SHELL:=/bin/bash

DOCKER:=docker
COMPOSE_DEV:=deployments/docker-compose.dev.yaml
COMPOSE_PROD:=deployments/docker-compose.yaml
COMPOSE_ENV:=deployments/compose.env
CONTAINERS:=db-media-equipment-tracker

define compose_file
	$(if $(findstring dev-,$(1)),$(COMPOSE_DEV),$(COMPOSE_PROD))
endef

%up:
	$(DOCKER) compose --env-file $(COMPOSE_ENV) -f $(call compose_file,$@) up -d $(CONTAINERS)

%upd:
	$(DOCKER) compose --env-file $(COMPOSE_ENV) -f $(call compose_file,$@) up -d --build $(CONTAINERS)

%upda:
	$(DOCKER) compose --env-file $(COMPOSE_ENV) -f $(call compose_file,$@) up --build $(CONTAINERS)

%down:
	$(DOCKER) compose --env-file $(COMPOSE_ENV) -f $(call compose_file,$@) down

