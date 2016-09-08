include ../build.mk

all: prepare
	@export builddir=$(shell pwd)/build; \
	echo $(SVNVERSION); \
	for pkg in $(PKG); do \
		pkg=$$(basename $${pkg} .spec); \
		echo "Building $${pkg}"; \
		paths=$$(rpm -q --specfile --qf "build/RPMS/%{arch}/$${pkg}-${SVNVERSION}-%{release}.%{arch}.rpm " $${pkg}.spec); \
		echo $${paths};\
		rpmbuild --define="_module_name ${PROJ}" --define="_version ${SVNVERSION}" --define="_topdir $${builddir}" -bb $${pkg}.spec && \
			pwd && \
			echo $${paths};\
			echo $${pkg};\
			mkdir -p backup-rpms ; cp $${paths} backup-rpms && \
			cp $${paths} ./ && \
			cp $${paths} ../ && \
			continue; \
		echo "Build failed for $${pkg}, abort."; \
		break; \
	done

# Make rpm building base and source tarball
prepare: clean-build clean-rpm
	echo "Creating source package to build/SOURCES/$(SRC)"
	mkdir -p $(TEMP_DIR)
	rm -f $(TEMP_DIR)/$(PROJ)
	ln -s $(ROOT_DIR) $(TEMP_DIR)/$(PROJ)
	$(MAKE) clean -C $(ROOT_DIR)
	tar zchf $(TEMP_DIR)/$(SRC) --exclude request_sample --exclude qmodule --exclude "*lib${PROJ}*.tar.gz" --exclude .git --exclude "*.rpm" -C $(TEMP_DIR) $(PROJ)
	mkdir -p build/{BUILD,RPMS,SRPMS,SPECS,SOURCES}
	mv $(TEMP_DIR)/$(SRC) build/SOURCES/

clean-build:
	rm -rf build

clean-rpm:
	rm -rf *rpm

clean: clean-build clean-rpm

.PHONY: pkg prepare clean clean-build clean-rpm

