import {Component, Inject, Input, OnDestroy} from '@angular/core';
import {Agent} from "../models/agent";
import {Apollo, gql} from "apollo-angular";
import {SubSink} from "subsink";
import {AlertService} from "../alert/alert.service";
import {MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialog as MatDialog} from "@angular/material/legacy-dialog";

@Component({
    selector: 'app-network-agents',
    template: `
        <table mat-table [dataSource]="agents">

            <ng-container matColumnDef="created">
                <th mat-header-cell *matHeaderCellDef> Created</th>
                <td mat-cell *matCellDef="let agent"> {{agent.created}} </td>
            </ng-container>
            <ng-container matColumnDef="name">
                <th mat-header-cell *matHeaderCellDef> Name</th>
                <td mat-cell *matCellDef="let agent"> {{agent.name}} </td>
            </ng-container>
            <ng-container matColumnDef="groups">
                <th mat-header-cell *matHeaderCellDef> Groups</th>
                <td mat-cell *matCellDef="let agent"> {{agent.groups}} </td>
            </ng-container>
            <ng-container matColumnDef="assignedIP">
                <th mat-header-cell *matHeaderCellDef> Nebula IP</th>
                <td mat-cell *matCellDef="let agent"> {{agent.assignedIP}} </td>
            </ng-container>
            <ng-container matColumnDef="actions">
                <th mat-header-cell *matHeaderCellDef></th>
                <td mat-cell *matCellDef="let agent">
                    <button mat-mini-fab color="warn" (click)="revokeAgentAccess(agent)">
                        <mat-icon>delete</mat-icon>
                    </button>
                </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="agentDisplayedColumns"></tr>
            <tr mat-row *matRowDef="let user; columns: agentDisplayedColumns;"></tr>
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
export class NetworkAgentsComponent implements OnDestroy {
    private subs = new SubSink();
    agentDisplayedColumns: string[] = ['created', 'name', 'groups', 'assignedIP', 'actions'];
    @Input() agents: Agent[] = [];

    constructor(private apollo: Apollo, private dialog: MatDialog, private alert: AlertService) {

    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }

    revokeAgentAccess(agent: Agent) {
        let dialogRef = this.dialog.open(NetworkAgentRevokeDialog, {
            data: agent
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.subs.sink = this.apollo.mutate({
                    variables: {
                        fingerprint: agent.fingerprint
                    },
                    mutation: gql`mutation RevokeCertsForAgent($fingerprint: String!) {
                        revokeCertsForAgent(fingerprint: $fingerprint)
                    }`,
                    update: (cache) => {
                        const normalizedId = cache.identify({
                            __typename: 'Agent',
                            fingerprint: agent.fingerprint,
                        });

                        cache.evict({id: normalizedId});
                        cache.gc();
                    }
                }).subscribe(() => {
                }, error => {
                    this.alert.addAlert('danger', error.message);
                })
            }
        });
    }
}


@Component({
    selector: 'agent-revoke-dialog',
    template: `<h1 mat-dialog-title>Revoke Agent</h1>
    <mat-dialog-content class="mat-typography">
        Fingerprint: {{agent.fingerprint}} <br/>
        Name: {{agent.name}} <br/>
        Network: {{agent.networkName}} <br/>
        Nebula IP: {{agent.assignedIP}} <br/>
        Groups: {{agent.groups}}
    </mat-dialog-content>
    <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-button color="warn" [mat-dialog-close]="true">Revoke</button>
    </mat-dialog-actions>`
})
export class NetworkAgentRevokeDialog {
    constructor(@Inject(MAT_DIALOG_DATA) public agent: Agent) {
    }

}
