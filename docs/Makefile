THEME_VERSION := v0.8.2
THEME := hugo-geekdoc
BASEDIR := .
THEMEDIR := $(BASEDIR)/themes

.PHONY: doc
doc: doc-assets doc-build

.PHONY: doc-assets
doc-assets:
	mkdir -p $(THEMEDIR)/$(THEME)/ ; \
	curl -sSL "https://github.com/thegeeklab/$(THEME)/releases/download/${THEME_VERSION}/$(THEME).tar.gz" | tar -xz -C $(THEMEDIR)/$(THEME)/ --strip-components=1

.PHONY: doc-build
doc-build:
	cd $(BASEDIR); hugo

.PHONY: clean
clean:
	rm -rf $(THEMEDIR) && \
	rm -rf $(BASEDIR)/public

start:
	hugo server -D