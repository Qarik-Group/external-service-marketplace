External Service Marketplace API Documentation
==============================================

This document lays out the design of the ESM front-end API, which
will be used by clients to provision, deprovision, bind, and
unbind service instances to their applications.

## GET /esm/v1/clouds

Retrieve a list of all known clouds to which service instances can
be bound (given an application running on that cloud).

## GET /esm/v1/clouds/:cloud

Retrieve the details for a single cloud, from the configuration,
but without passing back any sensitive information.

## GET /esm/v1/services/catalog

Retrieve the combined catalog, made by aggregating all of the
backend service brokers into one.  Each service broker's service
set will be prefixed with the brokers prefix, according to the
configuration set at boot time.

    GET /esm/v1/services/catalog

    {
        "services": [
          {
            "id": "tweed1/postgres",
            "name": "PostgreSQL Service",
            "description": "An RDBMS for a more civilized time",
            "tags": ["pg", "psql", "rad"],
            "bindable": true,
            // etc.

            "plans": [
            ]
          }
        ]
    }

## GET /esm/v1/services/:namespace/:service/:plan/instances

Retrieve all instances of the plan (requires ESM operator
privilege!)

## POST /esm/v1/services/:namespace/:service/:plan/instances

Creates an instance of the specified (namespaced) service and
plan, and returns information about that instance.

    POST /esm/v1/services/tweed1/postgres/small
    {}

Responses will be:

    {
      "instance": "foo",
    }

## GET /esm/v1/services/:namespace/:service/:plan/instances/:instance

Retrieve the details for a single instance.

## DELETE /esm/v1/services/:namespace/:service/:plan/instances

Delete an instance, and all of its bindings.

## GET /esm/v1/services/:namespace/:service/:plan/instances/:instance/bindings

Retrieve all bindings for a single instance.

## POST /esm/v1/services/:namespace/:service/:plan/instances/:instance/bindings

Bind (by creating a binding)

    POST /esm/v1/services/tweed1/pg/large/instances/d515bb35/bindings
    {
      "cloud": "c9f4a595", // an internal cloud id, from GET /esm/v1/clouds
      "namespace": "org/space",
      "app": "app-name",
    }

## GET /esm/v1/services/:namespace/:service/:plan/instances/:instance/bindings/:binding

Retrieve the details for a single binding of the given instance.

## DELETE /esm/v1/services/:namespace/:service/:plan/instances/:instance/bindings/:binding

Unbind (delete a binding).
