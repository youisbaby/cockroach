// Code generated by "stringer"; DO NOT EDIT.

package clusterversion

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[V21_2-0]
	_ = x[Start22_1-1]
	_ = x[ValidateGrantOption-2]
	_ = x[PebbleFormatBlockPropertyCollector-3]
	_ = x[ProbeRequest-4]
	_ = x[SelectRPCsTakeTracingInfoInband-5]
	_ = x[PreSeedTenantSpanConfigs-6]
	_ = x[SeedTenantSpanConfigs-7]
	_ = x[PublicSchemasWithDescriptors-8]
	_ = x[EnsureSpanConfigReconciliation-9]
	_ = x[EnsureSpanConfigSubscription-10]
	_ = x[EnableSpanConfigStore-11]
	_ = x[ScanWholeRows-12]
	_ = x[SCRAMAuthentication-13]
	_ = x[UnsafeLossOfQuorumRecoveryRangeLog-14]
	_ = x[AlterSystemProtectedTimestampAddColumn-15]
	_ = x[EnableProtectedTimestampsForTenant-16]
	_ = x[DeleteCommentsWithDroppedIndexes-17]
	_ = x[RemoveIncompatibleDatabasePrivileges-18]
	_ = x[AddRaftAppliedIndexTermMigration-19]
	_ = x[PostAddRaftAppliedIndexTermMigration-20]
	_ = x[DontProposeWriteTimestampForLeaseTransfers-21]
	_ = x[EnablePebbleFormatVersionBlockProperties-22]
	_ = x[DisableSystemConfigGossipTrigger-23]
	_ = x[MVCCIndexBackfiller-24]
	_ = x[EnableLeaseHolderRemoval-25]
	_ = x[BackupResolutionInJob-26]
	_ = x[LooselyCoupledRaftLogTruncation-27]
	_ = x[ChangefeedIdleness-28]
	_ = x[BackupDoesNotOverwriteLatestAndCheckpoint-29]
	_ = x[EnableDeclarativeSchemaChanger-30]
	_ = x[RowLevelTTL-31]
	_ = x[PebbleFormatSplitUserKeysMarked-32]
	_ = x[IncrementalBackupSubdir-33]
	_ = x[DateStyleIntervalStyleCastRewrite-34]
	_ = x[EnableNewStoreRebalancer-35]
	_ = x[ClusterLocksVirtualTable-36]
	_ = x[AutoStatsTableSettings-37]
	_ = x[ForecastStats-38]
	_ = x[SuperRegions-39]
	_ = x[EnableNewChangefeedOptions-40]
	_ = x[SpanCountTable-41]
	_ = x[PreSeedSpanCountTable-42]
	_ = x[SeedSpanCountTable-43]
	_ = x[V22_1-44]
	_ = x[Start22_2-45]
	_ = x[LocalTimestamps-46]
	_ = x[EnsurePebbleFormatVersionRangeKeys-47]
	_ = x[EnablePebbleFormatVersionRangeKeys-48]
	_ = x[TrigramInvertedIndexes-49]
	_ = x[RemoveGrantPrivilege-50]
	_ = x[MVCCRangeTombstones-51]
	_ = x[UpgradeSequenceToBeReferencedByID-52]
	_ = x[SampledStmtDiagReqs-53]
	_ = x[AddSSTableTombstones-54]
	_ = x[SystemPrivilegesTable-55]
	_ = x[EnablePredicateProjectionChangefeed-56]
	_ = x[AlterSystemSQLInstancesAddLocality-57]
}

const _Key_name = "V21_2Start22_1ValidateGrantOptionPebbleFormatBlockPropertyCollectorProbeRequestSelectRPCsTakeTracingInfoInbandPreSeedTenantSpanConfigsSeedTenantSpanConfigsPublicSchemasWithDescriptorsEnsureSpanConfigReconciliationEnsureSpanConfigSubscriptionEnableSpanConfigStoreScanWholeRowsSCRAMAuthenticationUnsafeLossOfQuorumRecoveryRangeLogAlterSystemProtectedTimestampAddColumnEnableProtectedTimestampsForTenantDeleteCommentsWithDroppedIndexesRemoveIncompatibleDatabasePrivilegesAddRaftAppliedIndexTermMigrationPostAddRaftAppliedIndexTermMigrationDontProposeWriteTimestampForLeaseTransfersEnablePebbleFormatVersionBlockPropertiesDisableSystemConfigGossipTriggerMVCCIndexBackfillerEnableLeaseHolderRemovalBackupResolutionInJobLooselyCoupledRaftLogTruncationChangefeedIdlenessBackupDoesNotOverwriteLatestAndCheckpointEnableDeclarativeSchemaChangerRowLevelTTLPebbleFormatSplitUserKeysMarkedIncrementalBackupSubdirDateStyleIntervalStyleCastRewriteEnableNewStoreRebalancerClusterLocksVirtualTableAutoStatsTableSettingsForecastStatsSuperRegionsEnableNewChangefeedOptionsSpanCountTablePreSeedSpanCountTableSeedSpanCountTableV22_1Start22_2LocalTimestampsEnsurePebbleFormatVersionRangeKeysEnablePebbleFormatVersionRangeKeysTrigramInvertedIndexesRemoveGrantPrivilegeMVCCRangeTombstonesUpgradeSequenceToBeReferencedByIDSampledStmtDiagReqsAddSSTableTombstonesSystemPrivilegesTableEnablePredicateProjectionChangefeedAlterSystemSQLInstancesAddLocality"

var _Key_index = [...]uint16{0, 5, 14, 33, 67, 79, 110, 134, 155, 183, 213, 241, 262, 275, 294, 328, 366, 400, 432, 468, 500, 536, 578, 618, 650, 669, 693, 714, 745, 763, 804, 834, 845, 876, 899, 932, 956, 980, 1002, 1015, 1027, 1053, 1067, 1088, 1106, 1111, 1120, 1135, 1169, 1203, 1225, 1245, 1264, 1297, 1316, 1336, 1357, 1392, 1426}

func (i Key) String() string {
	if i < 0 || i >= Key(len(_Key_index)-1) {
		return "Key(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Key_name[_Key_index[i]:_Key_index[i+1]]
}
