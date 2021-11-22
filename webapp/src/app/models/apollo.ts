import {Network} from "./network";
import {User} from "./user";
import {Agent} from "./agent";

export interface ApolloResponse {
    getAgent: Agent;
    getNetworks: Network[];
    getNetwork: Network;
    getUsers: User[];
    approveEnrollmentRequest: Agent;
}