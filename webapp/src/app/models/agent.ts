import {gql} from "apollo-angular";

export interface Agent {
    fingerprint?: string;
    created?: string;
    networkName: string;
    groups: string[];
    assignedIP: string;
    issuedAt: string;
    expiresAt: string;
    name: string;
}

export const GET_AGENT_BY_CLIENT_FINGERPRINT = gql`query GetAgent($fingerprint: String!){
    getAgent(fingerprint: $fingerprint){
        created
        fingerprint
        assignedIP
        groups
        name
    }
}`;