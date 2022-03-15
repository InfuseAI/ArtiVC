THEME_VERSION := v0.27.4
THEME := hugo-geekdoc
BASEDIR := .
THEMEDIR := $(BASEDIR)/themes

build: doc-assets doc-build

start: doc-assets
	hugo server -D	

doc-assets:
	mkdir -p $(THEMEDIR)/$(THEME)/ ; \
	curl -sSL "https://github.com/thegeeklab/$(THEME)/releases/download/${THEME_VERSION}/$(THEME).tar.gz" | tar -xz -C $(THEMEDIR)/$(THEME)/ --strip-components=1


doc-commands:
	rm -rf content/en/commands/art.md
	rm -rf content/en/commands/art_*.md
	go run ../main.go docs
	cp -R ./generated_docs/ content/en/commands/
	rm -rf ./generated_docs/

doc-build:
	cd $(BASEDIR); hugo

clean:
	rm -rf $(THEMEDIR) && \
	rm -rf $(BASEDIR)/public
