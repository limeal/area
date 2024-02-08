FROM cirrusci/android-sdk:33

USER root

ENV FLUTTER_VERSION "3.7.6"

WORKDIR /build

ENV FLUTTER_HOME=/build

ENV PATH ${PATH}:${FLUTTER_HOME}/bin:${FLUTTER_HOME}/bin/cache/dart-sdk/bin

RUN git clone --depth 1 --branch ${FLUTTER_VERSION} https://github.com/flutter/flutter.git ${FLUTTER_HOME}

RUN flutter doctor --android-licenses \
    && flutter doctor \
    && chown -R root:root ${FLUTTER_HOME}

COPY . .
RUN yes | flutter build apk --no-sound-null-safety --release
CMD ["cp", "build/app/outputs/flutter-apk/app-release.apk", "/data/app-release.apk"]