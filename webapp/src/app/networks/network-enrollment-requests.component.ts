import {Component, Input, OnDestroy} from '@angular/core';
import {EnrollmentRequest, GET_NETWORK_BY_NAME} from "../models/network";
import {ApolloResponse} from "../models/apollo";
import {Apollo, gql} from "apollo-angular";
import {SubSink} from "subsink";
import {MatLegacyDialog as MatDialog} from "@angular/material/legacy-dialog";
import {AlertService} from "../alert/alert.service";
import {EnrollmentRequestApproveDialog} from "./enrollment-request-approve-dialog.component";

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
            <ng-container matColumnDef="groups">
                <th mat-header-cell *matHeaderCellDef> Groups</th>
                <td mat-cell *matCellDef="let er"> {{er.groups}} </td>
            </ng-container>
            <ng-container matColumnDef="requestedIP">
                <th mat-header-cell *matHeaderCellDef> Requested Nebula IP</th>
                <td mat-cell *matCellDef="let er"> {{er.requestedIP}} </td>
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
    displayedColumns: string[] = ['created', 'name', 'clientIP', 'groups', 'requestedIP', 'actions'];
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
                        fingerprint: er.fingerprint
                    },
                    mutation: gql`mutation ApproveEnrollmentRequest($fingerprint: String!) {
                        approveEnrollmentRequest(fingerprint: $fingerprint){
                            created
                            fingerprint
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
                fingerprint: er.fingerprint
            },
            mutation: gql`mutation DeleteEnrollmentRequest($fingerprint: String!) {
                deleteEnrollmentRequest(fingerprint: $fingerprint)
            }`,
            update: (cache) => {
                const normalizedId = cache.identify({
                    __typename: 'EnrollmentRequest',
                    fingerprint: er.fingerprint,
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


