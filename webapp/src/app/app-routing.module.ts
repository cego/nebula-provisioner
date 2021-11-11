import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {NetworksComponent} from "./networks/networks.component";
import {UsersComponent} from "./users/users.component";
import {NetworkComponent} from "./networks/network.component";

const routes: Routes = [
    {path: "networks", component: NetworksComponent},
    // {path: "networks/add", component: NetworkAddComponent}, // TODO Change route to not conflict with network name
    {path: "networks/:name", component: NetworkComponent},
    {path: "users", component: UsersComponent},
    {path: "**", redirectTo: "networks"}
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
