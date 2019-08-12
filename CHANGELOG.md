# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
- Remove MobileSecurityServicePodCount since it was considered duplicated [#157](https://github.com/aerogear/mobile-security-service-operator/pull/157)
- Fix the chart `Mobile Security Service Application - Uptime` on Grafana application dashboard [#164](https://github.com/aerogear/mobile-security-service-operator/pull/164)
- Fix Grafana application resources dashboard [#165](https://github.com/aerogear/mobile-security-service-operator/pull/165)
- Upgrade version of operator-sdk from 0.8.1 to 0.10.0 [#163](https://github.com/aerogear/mobile-security-service-operator/pull/163)

## [0.3.0] - 2019-07-26
- Fixed Prometheus Rules for MobileSecurityServicePodCount and MobileSecurityServiceDown [#151](https://github.com/aerogear/mobile-security-service-operator/pull/151)
- Improved and fixed sop in order to add steps to scale the pod [#151](https://github.com/aerogear/mobile-security-service-operator/pull/151)
- Added schema validation by OpenAPI [#149](https://github.com/aerogear/mobile-security-service-operator/pull/149)
- Fixed GVK group to use "aerogear.org"
- Updated Deployments version and group from extensions v1beta1 to apps v1
- Removed some unnecessary permissions from the operator ServiceAccount

## [0.2.0] - 2019-07-01
- Release of the operator which meets all criteria planned into the https://trello.com/b/mCiUFubz/security-service.

## [0.1.1] - [0.1.4]
- Improvements and fixes performed in the initial development phase

## [0.1.0] - 2019-06-12
- Initial release of operator that meets requirements to manage [MobileSecurityService](https://github.com/aerogear/mobile-security-service) and its database.
