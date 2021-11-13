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
    clientFingerprint?: string
    created?: string
    networkName?: string
    clientIP?: string
    name?: string
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
            clientFingerprint
            clientIP
            created
            name
            networkName
        }
        agents {
            created
            clientFingerprint
            assignedIP
            groups
            name
        }
    }
}`;