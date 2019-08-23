# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
- Change/Add Memory and CPU resources requested and limits for MSS, Oauth and Operator [#172](https://github.com/aerogear/mobile-security-service-operator/pull/172)

## [0.4.0] - 2019-08-13
- Remove MobileSecurityServicePodCount since it was considered duplicated [#157](https://github.com/aerogear/mobile-security-service-operator/pull/157)
- Fix the chart `Mobile Security Service Application - Uptime` on Grafana application dashboard [#164](https://github.com/aerogear/mobile-security-service-operator/pull/164)
- Fix Grafana application resources dashboard [#165](https://github.com/aerogear/mobile-security-service-operator/pull/165)
- Fix Monitor Operator metrics which were not working because the `service/mobile-security-service-operator` was not been created by the operator [#166](https://github.com/aerogear/mobile-security-service-operator/pull/166)
- Removed the creation of a config map created for each MobileSecurityServiceApp CR as it is not needed. [#167](https://github.com/aerogear/mobile-security-service-operator/pull/167)
- Upgrade the version of Mobile Security Service used by default from 0.1.0 to 0.2.2 [#168](https://github.com/aerogear/mobile-security-service-operator/pull/168)

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
