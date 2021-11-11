export interface User {
    id?: string;
    name?: string;
    email?: string;
    userApprove?: UserApprove;
    disabled?: boolean;
}

export interface UserApprove {
    approved?: boolean;
    approvedBy?: string;
    approvedAt?: string;
    approvedByUser?: User;
}