BuildEnv Docker versions
========================
These are the versions of objectboxio/buildenv-generator-ubuntu-* images.
Note: for each version, please indicate the primary motivation(s) that led to creating the version.

objectboxio/buildenv-generator-ubuntu:2025-10-16
------------------------------------------------
Update to the latest Go version for the upcoming 5.0 release.

```
=== Go versions (system installed first) ===
which: /usr/bin/go
go version go1.18.1 linux/amd64
which: /usr/local/go1.24.9/bin/go
go version go1.24.9 linux/amd64
which: /usr/local/go1.25/bin/go
go version go1.25.3 linux/amd64

=== Compiler, libs & Tools ===
GNU C Library (Ubuntu GLIBC 2.35-0ubuntu3.11) stable release version 2.35.
GLIBCXX_3.4
GLIBCXX_3.4.30
/usr/bin/gcc
gcc (Ubuntu 11.4.0-1ubuntu1~22.04.2) 11.4.0
/usr/bin/ld
GNU ld (GNU Binutils for Ubuntu) 2.38
cmake version 3.22.1
GNU Make 4.3
ccache version 4.5.1

REPOSITORY                              TAG          IMAGE ID       CREATED        SIZE
objectboxio/buildenv-generator-ubuntu   2025-10-16   538e75e5d488   1 second ago   1.34GB
```

objectboxio/buildenv-generator-ubuntu:2024-02-26
------------------------------------------------

| Version | Location (GOROOT)   | Notes                                            |
|---------|---------------------|--------------------------------------------------|
| 1.18.1  | /usr/lib/go-1.18    | (Ubuntu 22.04 package)                           |
| 1.22.0  | /usr/local/go1.22   | Manual installed version (default first on PATH) |
| 1.19.13 | /root/sdk/go1.19.13 | Additional version installed via go install      |
