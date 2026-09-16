#!/bin/bash

rm generated/*

kustomize build --load-restrictor LoadRestrictionsNone base/bootstrap > ./generated/bootstrap-base.yaml
kustomize build --load-restrictor LoadRestrictionsNone base/controlplane > ./generated/controlplane-base.yaml
kustomize build --load-restrictor LoadRestrictionsNone global/bootstrap > ./generated/bootstrap-global.yaml
kustomize build --load-restrictor LoadRestrictionsNone global/controlplane > ./generated/controlplane-global.yaml
