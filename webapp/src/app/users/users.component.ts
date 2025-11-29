import {Component, Inject, OnDestroy, OnInit} from '@angular/core';
import {of as observableOf} from "rxjs";
import {catchError, map} from "rxjs/operators";
import {Apollo, gql} from "apollo-angular";
import {ApolloResponse} from "../models/apollo";
import {SubSink} from "subsink";
import {User} from "../models/user";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {AlertService} from "../alert/alert.service";

@Component({
    selector: 'app-users',
    template: `
        <div class="mat-elevation-z8">
            <div class="loading-shade" *ngIf="isLoadingResults">
                <mat-spinner *ngIf="isLoadingResults"></mat-spinner>
            </div>

            <table mat-table [dataSource]="data">

                <ng-container matColumnDef="name">
                    <th mat-header-cell *matHeaderCellDef> Name</th>
                    <td mat-cell *matCellDef="let user"> {{user.name}} </td>
                </ng-container>
                <ng-container matColumnDef="email">
                    <th mat-header-cell *matHeaderCellDef> E-mail</th>
                    <td mat-cell *matCellDef="let user"> {{user.email}} </td>
                </ng-container>
                <ng-container matColumnDef="approve">
                    <th mat-header-cell *matHeaderCellDef></th>
                    <td mat-cell *matCellDef="let user">
                        <button mat-raised-button (click)="approveUserDialog(user)"
                                *ngIf="!user.userApprove?.approved || user.disabled">
                            Approve
                        </button>
                        <span *ngIf="user.userApprove?.approved && !user.disabled">Approved by: {{user.userApprove?.approvedByUser?.name}}</span>
                    </td>
                </ng-container>
                <ng-container matColumnDef="actions">
                    <th mat-header-cell *matHeaderCellDef></th>
                    <td mat-cell *matCellDef="let user">
                        <button mat-mini-fab color="warn" *ngIf="!user.disabled" (click)="deleteUserDialog(user)">
                            <mat-icon>delete</mat-icon>
                        </button>
                    </td>
                </ng-container>


                <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
                <tr mat-row *matRowDef="let user; columns: displayedColumns;"></tr>
            </table>
        </div>
    `,
    styles: [`
      table {
        width: 100%;
      }

      .action-buttons {
        padding: 1em;
        margin-bottom: 1em;
      }

      .mat-column-approve {
        width: 20em;
        text-align: left;
      }

      .mat-column-actions {
        padding-top: 0.5em;
        padding-bottom: 0.5em;
        width: 1em;
        text-align: right;
      }

    `],
    standalone: false
})
export class UsersComponent implements OnInit, OnDestroy {
    private subs = new SubSink();
    displayedColumns: string[] = ['name', 'email', 'approve', 'actions'];
    data: User[] = [];
    expandedUser: User | null = null;

    isLoadingResults = false;

    constructor(private apollo: Apollo, private dialog: MatDialog, private alert: AlertService) {
    }

    ngOnInit(): void {
        this.isLoadingResults = true;

        this.subs.sink = this.apollo.watchQuery<ApolloResponse>({
            query: gql`
                {
                    getUsers {
                        id
                        name
                        email
                        disabled
                        userApprove {
                            approved
                            approvedBy
                            approvedAt
                            approvedByUser {
                                name
                            }
                        }
                    }
                }`,
        }).valueChanges.pipe(map(res => {
                this.isLoadingResults = false;
                return res.data.getUsers;
            }),
            catchError(() => observableOf([])))
        .subscribe(data => {
            this.data = data;
        }, error => {
            this.alert.addAlert('danger', error.message);
        });

    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }

    approveUserDialog(user: User) {
        let dialogRef = this.dialog.open(UserApproveDialog, {
            data: user
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.subs.sink = this.apollo.mutate({
                    variables: {
                        userId: user.id
                    },
                    mutation: gql`mutation ApproveUser($userId: String!) {
                        approveUser(userId: $userId){
                            id
                            disabled
                            userApprove {
                                approved
                                approvedBy
                                approvedAt
                                approvedByUser {
                                    id
                                    name
                                }
                            }
                        }
                    }`
                }).subscribe(_ => {
                }, error => {
                    this.alert.addAlert('danger', error.message);
                });
            }
        });
    }

    deleteUserDialog(user: User) {
        let dialogRef = this.dialog.open(UserDeleteDialog, {
            data: user
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.subs.sink = this.apollo.mutate({
                    variables: {
                        userId: user.id
                    },
                    mutation: gql`mutation disableUser($userId: String!) {
                        disableUser(userId: $userId){
                            id
                            disabled
                            userApprove {
                                approved
                                approvedBy
                                approvedAt
                                approvedByUser {
                                    id
                                    name
                                }
                            }
                        }
                    }`
                }).subscribe(_ => {
                }, error => {
                    this.alert.addAlert('danger', error.message);
                })
            }
        });
    }
}

@Component({
    selector: 'user-approve-dialog',
    template: `<h1 mat-dialog-title>Approve User: {{user.email}}</h1>
    <mat-dialog-content class="mat-typography">
        ID: {{user.id}}<br/>
        Name: {{user.name}}<br/>
        E-mail: {{user.email}}
    </mat-dialog-content>
    <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-button color="warn" [mat-dialog-close]="true">Approve</button>
    </mat-dialog-actions>
    `,
    standalone: false
})
export class UserApproveDialog {
    constructor(@Inject(MAT_DIALOG_DATA) public user: User) {
    }
}

@Component({
    selector: 'user-delete-dialog',
    template: `<h1 mat-dialog-title>Disable User: {{user.email}}</h1>
    <mat-dialog-content class="mat-typography">
        ID: {{user.id}}<br/>
        Name: {{user.name}}<br/>
        E-mail: {{user.email}}
    </mat-dialog-content>
    <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-button color="warn" [mat-dialog-close]="true">Disable</button>
    </mat-dialog-actions>
    `,
    standalone: false
})
export class UserDeleteDialog {
    constructor(@Inject(MAT_DIALOG_DATA) public user: User) {
    }
}