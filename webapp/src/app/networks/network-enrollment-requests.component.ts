import {Component, Inject, Input, OnDestroy} from '@angular/core';
import {EnrollmentRequest, GET_NETWORK_BY_NAME} from "../models/network";
import {ApolloResponse} from "../models/apollo";
import {Apollo, gql} from "apollo-angular";
import {SubSink} from "subsink";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {AlertService} from "../alert/alert.service";

@Component({
    selector: 'app-network-enrollment-requests',
    template: `
        <table mat-table [dataSource]="enrollmentRequests">

            <ng-container matColumnDef="created">
                <th mat-header-cell *matHeaderCellDef> Created</th>
                <td mat-cell *matCellDef="let er"> {{er.created}} </td>
            </ng-container>
            <ng-container matColumnDef="name">
                <th mat-header-cell *matHeaderCellDef> Name</th>
                <td mat-cell *matCellDef="let er"> {{er.name}} </td>
            </ng-container>
            <ng-container matColumnDef="clientIP">
                <th mat-header-cell *matHeaderCellDef> Requested from IP</th>
                <td mat-cell *matCellDef="let er"> {{er.clientIP}} </td>
            </ng-container>
            <ng-container matColumnDef="actions">
                <th mat-header-cell *matHeaderCellDef></th>
                <td mat-cell *matCellDef="let er">
                    <button mat-mini-fab color="primary" (click)="approveEnrollmentRequestDialog(er)">
                        <mat-icon>done</mat-icon>
                    </button>
                    <button mat-mini-fab color="warn" (click)="deleteEnrollmentRequest(er)">
                        <mat-icon>delete</mat-icon>
                    </button>
                </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let user; columns: displayedColumns;"></tr>
        </table>
    `,
    styles: [`
      table {
        width: 100%;
      }

      mat-form-field {
        width: 100%;
        margin: 1em;
      }

      .mat-column-actions {
        padding-top: 0.5em;
        padding-bottom: 0.5em;
        width: 10em;
        text-align: right;
      }
    `],
})
export class NetworkEnrollmentRequestsComponent implements OnDestroy {
    private subs = new SubSink();
    displayedColumns: string[] = ['created', 'name', 'clientIP', 'actions'];
    @Input() enrollmentRequests: EnrollmentRequest[] = [];

    constructor(private apollo: Apollo, private dialog: MatDialog, private alert: AlertService) {

    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }


    approveEnrollmentRequestDialog(er: EnrollmentRequest) {
        let dialogRef = this.dialog.open(EnrollmentRequestApproveDialog, {
            data: er
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.subs.sink = this.apollo.mutate<ApolloResponse>({
                    variables: {
                        clientFingerprint: er.clientFingerprint
                    },
                    mutation: gql`mutation ApproveEnrollmentRequest($clientFingerprint: String!) {
                        approveEnrollmentRequest(clientFingerprint: $clientFingerprint){
                            created
                            clientFingerprint
                            assignedIP
                            groups
                            name
                        }
                    }`,
                    refetchQueries: [
                        GET_NETWORK_BY_NAME
                    ]
                }).subscribe(() => {
                }, error => {
                    this.alert.addAlert('danger', error.message);
                })
            }
        });
    }

    deleteEnrollmentRequest(er: EnrollmentRequest) {
        this.subs.sink = this.apollo.mutate({
            variables: {
                clientFingerprint: er.clientFingerprint
            },
            mutation: gql`mutation DeleteEnrollmentRequest($clientFingerprint: String!) {
                deleteEnrollmentRequest(clientFingerprint: $clientFingerprint)
            }`,
            update: (cache) => {
                const normalizedId = cache.identify({
                    __typename: 'EnrollmentRequest',
                    clientFingerprint: er.clientFingerprint,
                });

                cache.evict({id: normalizedId});
                cache.gc();
            }
        }).subscribe(() => {
        }, error => {
            this.alert.addAlert('danger', error.message);
        })
    }
}


@Component({
    selector: 'enrollment-request-approve-dialog',
    template: `<h1 mat-dialog-title>Approve Enrollment Request</h1>
    <mat-dialog-content class="mat-typography">
        Created: {{er.created}}<br/>
        Client fingerprint: {{er.clientFingerprint}}<br/>
        Name: {{er.name}}<br/>
        Client IP: {{er.clientIP}}<br/>
        Network Name: {{er.networkName}}<br/>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-button color="warn" [mat-dialog-close]="true">Approve</button>
    </mat-dialog-actions>
    `
})
export class EnrollmentRequestApproveDialog {
    constructor(@Inject(MAT_DIALOG_DATA) public er: EnrollmentRequest) {
    }
}
