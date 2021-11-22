import {Component, Input} from '@angular/core';
import {Agent} from "../models/agent";

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
export class NetworkAgentsComponent {
    agentDisplayedColumns: string[] = ['created', 'name', 'groups', 'assignedIP'];
    @Input() agents: Agent[] = [];

    constructor() {

    }
}