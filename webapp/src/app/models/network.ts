import {Agent} from "./agent";
import {gql} from "apollo-angular";

export interface Network {
    name?: string;
    ipPools?: string[];
    duration?: string;
    groups?: string[];
    ips?: string[];
    subnets?: string[];
    agents: Agent[];
    enrollmentToken?: string;
    enrollmentRequests?: EnrollmentRequest[]
}

export interface EnrollmentRequest {
    fingerprint?: string;
    created?: string;
    networkName?: string;
    clientIP?: string;
    name?: string;
    requestedIP: string;
    groups?: string[];
}

export const GET_NETWORK_BY_NAME = gql`query GetNetwork($name: String!){
    getNetwork(name: $name){
        name
        groups
        ipPools
        ips
        subnets
        enrollmentToken
        enrollmentRequests {
            fingerprint
            clientIP
            created
            name
            networkName
            requestedIP
            groups
        }
        agents {
            created
            fingerprint
            assignedIP
            groups
            name
        }
    }
}`;