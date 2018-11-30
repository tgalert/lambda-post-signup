const k8s = require('@kubernetes/client-node');

class KubeCtl {

    constructor() {
        const kc = new k8s.KubeConfig();
        kc.loadFromDefault();
        this.appsV1Api = kc.makeApiClient(k8s.Apps_v1Api);
    }

    createDeployment(name, rabbitUser, rabbitPw, rabbitVhost) {
        const deployment = Object.assign({}, deploymentSlug);
        setDeploymentName(deployment, name);
        setDeploymentAmqpUri(deployment, getAmqpUri(rabbitUser, rabbitPw, rabbitVhost));
        return this.appsV1Api.createNamespacedDeployment('default', deployment, true);
    }
}

module.exports = KubeCtl;

function setDeploymentAmqpUri(deployment, amqpUri) {
    const env = deployment.spec.template.spec.containers[0].env;
    for (const e of env) {
        if (e.name === 'AMQP_URI') e.value = amqpUri;
    }
}

function setDeploymentName(deployment, name) {
    deployment.metadata.name = name;
}

function getAmqpUri(user, password, vhost) {
    return `amqp://${user}:${password}@rabbitmq-service/${vhost}`;
}

// TODO: read secrets from secrets file
// TODO: move Docker image specification to config file
const deploymentSlug = {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {name: '[[REPLACE]]'},
    spec: {
        replicas: 1,
        selector: {matchLabels: {type: 'core'}},
        template: {
            metadata: {labels: {type: 'core'}},
            spec: {
                containers: [
                    {
                        name: 'core',
                        image: 'weibeld/tg-monitor:core-0.0.1',
                        env: [
                            {name: 'AMQP_URI', value: '[[REPLACE]]'},
                            {name: 'TG_API_ID', value: '208236'},
                            {name: 'TG_API_HASH', value: 'eaef680343e6b52e7011e8e9e442be86'},
                            {name: 'MAILGUN_API_KEY', value: 'key-524417edbb4ebd24ffc1ee9201dce78a'},
                            {name: 'MAILGUN_SENDING_ADDRESS', value: 'tg-monitor@quantumsense.ai'},
                            {name: 'MAILGUN_DOMAIN', value: 'quantumsense.ai'}
                        ]
                    }
                ]
            }
        }
    }
};


/* To read an object from a YAML file, use the following:
 *
 *   const yaml = require('js-yaml');
 *   const fs = require('fs');
 *   const deployment = yaml.safeLoad(fs.readFileSync('file.yml', 'utf8'));
 *
 * If there are multiple objects in the YAML file separated by '---', use:
 *
 *   const objects = yaml.safeLoadAll(fs.readFileSync('file.yml', 'utf8'));
 *
 * In this case, 'objects' is an array of objects.
 */