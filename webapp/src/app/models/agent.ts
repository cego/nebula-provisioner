export interface Agent {
    clientFingerprint?: string;
    created?: string;
    networkName: string;
    groups: string[];
    assignedIP: string;
    issuedAt: string;
    expiresAt: string;
    name: string;
}