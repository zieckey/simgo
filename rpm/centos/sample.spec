%define debug_package %{nil}
%define __os_install_post %{nil}

%define approot /home/s
%define binpath  %{approot}/%{_module_name}/bin
%define sbinpath %{approot}/%{_module_name}/sbin
%define logspath %{approot}/%{_module_name}/logs
%define confpath %{approot}/%{_module_name}/conf
%define datapath %{approot}/data/%{_module_name}

summary: The %{_module_name} project to check web page safe
name: %{_module_name}
version: %{_version}
release: 1%{?dist}
license: Commercial
vendor: Qihoo <http://www.360.cn>
group: Development/Libraries
source: %{_module_name}.tar.gz
requires: libqlog >= 2.1.1
requires: libcloudcom >= 1.5.1
requires: dconf_reload >= 1.2.1
buildrequires: libqlog-devel >= 2.1.1
buildrequires: libcloudcom-devel >= 1.5.1
buildroot: %{_tmppath}/%{name}-%{version}-%{release}-%(%{__id_u} -n)
Autoreq: no

%description
The %{_module_name} project to check web page safe

%prep
%setup -q -n %{_module_name}

%build
# your package build steps
make clean
make 

%install
rm -rf %{buildroot}
# your package install steps
# the compiled files dir: %{_builddir}/<package_source_dir> or $RPM_BUILD_DIR/<package_source_dir>
# the dest root dir: %{buildroot} or $RPM_BUILD_ROOT
mkdir -p %{buildroot}/%{binpath}
mkdir -p %{buildroot}/%{sbinpath}
mkdir -p %{buildroot}/%{logspath}
mkdir -p %{buildroot}/%{logspath}/framework
mkdir -p %{buildroot}/%{logspath}/access
mkdir -p %{buildroot}/%{logspath}/run
mkdir -p %{buildroot}/%{logspath}/stat
mkdir -p %{buildroot}/%{logspath}/check
mkdir -p %{buildroot}/%{confpath}
mkdir -p %{buildroot}/%{confpath}/common
mkdir -p %{buildroot}/%{datapath}

echo topdir: %{_topdir}
echo version: %{_version}
echo module_name: %{_module_name}
echo approot: %{approot}
echo buildroot: %{buildroot}

pushd %{_builddir}/%{_module_name}/bin
cp srvctl     %{buildroot}/%{binpath}

cp init.sh    %{buildroot}/%{binpath}
cp install.sh %{buildroot}/%{binpath}
popd

pushd %{_builddir}/%{_module_name}/sbin
cp wpe %{buildroot}/%{sbinpath}
popd

pushd %{_builddir}/%{_module_name}/conf
cp *.bin  %{buildroot}/%{confpath}
cp *.ini  %{buildroot}/%{confpath}
cp *.conf %{buildroot}/%{confpath}

cp status.html          %{buildroot}/%{confpath}
cp status.html.ok       %{buildroot}/%{confpath}/common
cp status.html.maintain %{buildroot}/%{confpath}/common
popd

# override
pushd %{_builddir}/%{_module_name}/conf/latest
cp *.conf %{buildroot}/%{confpath}
popd

%files
%defattr(-,cloud,cloud)
# list your package files here
# the list of the macros:
#   _prefix           /usr
#   _exec_prefix      %{_prefix}
#   _bindir           %{_exec_prefix}/bin
#   _libdir           %{_exec_prefix}/%{_lib}
#   _libexecdir       %{_exec_prefix}/libexec
#   _sbindir          %{_exec_prefix}/sbin
#   _includedir       %{_prefix}/include
#   _datadir          %{_prefix}/share
#   _sharedstatedir   %{_prefix}/com
#   _sysconfdir       /etc
#   _initrddir        %{_sysconfdir}/rc.d/init.d
#   _var              /var
%dir %{approot}
%dir %{binpath}
%dir %{sbinpath}
%dir %{confpath}
%dir %{confpath}/common
%dir %{datapath}
%dir %{logspath}
%dir %{logspath}/framework
%dir %{logspath}/access
%dir %{logspath}/stat
%dir %{logspath}/check
%dir %{logspath}/run

%attr(755,cloud,cloud) %{binpath}/srvctl
%attr(755,cloud,cloud) %{binpath}/init.sh
%attr(755,cloud,cloud) %{binpath}/install.sh

%attr(644,cloud,cloud) %config(noreplace,missingok) %{confpath}/app.conf
%attr(644,cloud,cloud) %config(noreplace,missingok) %{confpath}/qlog.conf

%{sbinpath}/wpe
%{confpath}/ngx_wpe_location.conf

%{confpath}/dconf_wpe.ini
%{confpath}/danger_level_file.ini

%{confpath}/asymmetric_keys.bin
%{confpath}/symmetric_keys.bin

%{confpath}/status.html
%{confpath}/common/status.html.ok
%{confpath}/common/status.html.maintain


%pre
# pre-install scripts

%post
# post-install scripts
%{binpath}/install.sh

%preun
# pre-uninstall scripts

%postun
# post-uninstall scripts

%clean
rm -rf %{buildroot}
# your package build clean up steps here

%changelog
# list your change log here

