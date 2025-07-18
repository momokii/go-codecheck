// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {docker} from '../models';
import {models} from '../models';
import {main} from '../models';
import {parser} from '../models';
import {semgrep} from '../models';

export function CheckDockerImagesIsAvailable(arg1:string,arg2:string):Promise<docker.DockerImage>;

export function CheckDockerIsAvailable():Promise<string>;

export function CheckIfFolderOrFIleExists(arg1:string):Promise<boolean>;

export function CreateNewRepo(arg1:models.RepositoryCreate):Promise<void>;

export function CreateNewScan(arg1:number,arg2:models.ScanCreate):Promise<void>;

export function DeleteRepo(arg1:number,arg2:number):Promise<void>;

export function DeleteScan(arg1:number,arg2:number):Promise<void>;

export function GetAndValidateUserByToken(arg1:string):Promise<models.User>;

export function GetRepoById(arg1:number,arg2:number):Promise<models.Repository>;

export function GetRepoDatas(arg1:number,arg2:number,arg3:number,arg4:string,arg5:boolean):Promise<main.RepoDataPaginationFE>;

export function GetScanById(arg1:number,arg2:number):Promise<models.ScanFull>;

export function GetScanDatas(arg1:number,arg2:number,arg3:number,arg4:string,arg5:boolean):Promise<main.ScanDataPaginationFE>;

export function GetScanDetail(arg1:string):Promise<parser.SemgrepReport>;

export function GetSemgrepReportData():Promise<parser.SemgrepReport>;

export function InitAndPrepareFolderScanSemgrep(arg1:string):Promise<void>;

export function Login(arg1:models.UserLogin):Promise<string>;

export function Logout(arg1:number):Promise<void>;

export function RunSemgrepScan():Promise<semgrep.ScanResult>;

export function UpdateRepo(arg1:number,arg2:models.RepositoryUpdate):Promise<void>;

export function UpdateUserPassword(arg1:models.UserUpdate,arg2:boolean):Promise<void>;

export function UpdateUserUsername(arg1:models.UserUpdate):Promise<void>;
