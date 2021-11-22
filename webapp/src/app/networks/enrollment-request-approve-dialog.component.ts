import {AfterViewInit, Component, Inject, OnDestroy} from "@angular/core";
import {MAT_DIALOG_DATA} from "@angular/material/dialog";
import {EnrollmentRequest} from "../models/network";
import {Agent, GET_AGENT_BY_CLIENT_FINGERPRINT} from "../models/agent";
import {Apollo} from "apollo-angular";
import {SubSink} from "subsink";
import {ApolloResponse} from "../models/apollo";
import {AlertService} from "../alert/alert.service";

@Component({
    selector: 'enrollment-request-approve-dialog',
    template: `<h1 mat-dialog-title>Approve Enrollment Request</h1>
    <mat-dialog-content class="mat-typography">

        <table mat-table [dataSource]="datasource">
            <ng-container matColumnDef="name">
                <th mat-header-cell *matHeaderCellDef></th>
                <td mat-cell *matCellDef="let element"> {{element.name}} </td>
            </ng-container>
            <ng-container matColumnDef="new">
                <th mat-header-cell *matHeaderCellDef> New</th>
                <td mat-cell *matCellDef="let element"> {{element.new}} </td>
            </ng-container>
            <ng-container matColumnDef="old">
                <th mat-header-cell *matHeaderCellDef> Current</th>
                <td mat-cell *matCellDef="let element"> {{element.old}} </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
        </table>

        <hr>
        Fingerprint: {{er.fingerprint}}
    </mat-dialog-content>
    <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-button color="warn" [mat-dialog-close]="true">Approve</button>
    </mat-dialog-actions>
    `,
    styles: [`
      table {
        width: 100%;
      }`]
})
export class EnrollmentRequestApproveDialog implements AfterViewInit, OnDestroy {
    private subs = new SubSink();
    existingAgent!: Agent
    displayedColumns: string[] = [];
    datasource: any[] = [];

    constructor(@Inject(MAT_DIALOG_DATA) public er: EnrollmentRequest,
                private apollo: Apollo,
                private alert: AlertService) {
    }

    ngAfterViewInit(): void {
        this.renderRows();
        this.displayedColumns = ['name', 'new'];
        if (this.er.fingerprint != null) {
            this.displayedColumns = ['name', 'old', 'new'];
            this.getAgent(this.er.fingerprint);
        }
    }

    renderRows() {
        this.datasource = [
            {
                name: 'Created',
                old: this.existingAgent?.created,
                new: this.er.created
            },
            {
                name: 'Agent Name',
                old: this.existingAgent?.name,
                new: this.er.name
            },
            {
                name: 'Nebula IP',
                old: this.existingAgent?.assignedIP,
                new: this.er.requestedIP
            },
            {
                name: 'Groups',
                old: this.existingAgent?.groups,
                new: this.er.groups
            }
        ];
    }

    getAgent(fingerprint: string) {
        this.subs.sink = this.apollo.query<ApolloResponse>({
            variables: {
                fingerprint: fingerprint
            },
            query: GET_AGENT_BY_CLIENT_FINGERPRINT,
        }).subscribe(data => {
            this.existingAgent = data.data.getAgent;
            this.renderRows();
        }, error => {
            this.alert.addAlert('danger', error.message);
        });
    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }
}