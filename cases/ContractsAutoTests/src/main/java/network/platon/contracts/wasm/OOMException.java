package network.platon.contracts.wasm;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.abi.datatypes.WasmFunction;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.WasmContract;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.WasmFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with platon-web3j version 0.13.1.5.
 */
public class OOMException extends WasmContract {
    private static String BINARY_0 = "0x0061736d0100000001550c60017f0060000060017f017f60027f7f0060027f7f017f60047f7f7f7f006000017f60047f7f7f7f017f60087f7f7f7f7f7f7f7f017f60097f7f7f7f7f7f7f7f7f017f60087f7f7e7f7e7f7f7f017f60017f017e025c0403656e760c706c61746f6e5f70616e6963000103656e760c706c61746f6e5f6465627567000303656e7617706c61746f6e5f6765745f696e7075745f6c656e677468000603656e7610706c61746f6e5f6765745f696e707574000003181701000000020107030400020b030104010205080a0902020405017001040405030100020608017f01419089040b073904066d656d6f72790200115f5f7761736d5f63616c6c5f63746f727300040f5f5f66756e63735f6f6e5f65786974001106696e766f6b6500090909010041010b031505060abf3b17040010130b0300010b5801027f230041206b2201240041002100034020004180304604402001200228020036020020011007200141206a240005418080c00010082202200036020020012000360210200141106a1007200041016a21000c010b0b0b8c19030f7f027e037c230041c0086b220424002004200036020c4180082105200441800836029c082004419f086a210e0340410020036b21060240034020052d00002201450d0120014125470440200341ff074d0440200441106a20036a20013a00000b200341016a21032004200541016a220536029c082006417f6a21060c010b0b2004200541016a220536029c08410021010340024002400240024020052c0000220841556a220241054b0440200841606a220241034b0d0102400240200241016b0e03030301000b2004200541016a220536029c08200141087221010c060b2004200541016a220536029c08200141107221010c050b200241016b0e050002000003010b0240200841506a41ff017141094d04402004419c086a1014210a200428029c0821050c010b4100210a2008412a470d00200028020021022004200541016a220536029c082001410272200120024100481b210120022002411f7522086a200873210a200041046a21000b41002108024020052d0000412e470d002004200541016a220736029c08200141800872210120052d0001220241506a41ff017141094d04402004419c086a10142108200428029c0821050c010b200241ff0171412a470440200721050c010b200028020021022004200541026a220536029c0820024100200241004a1b2108200041046a21000b0240024020052c000041987f6a411f77220241094b0d000240024002400240200241016b0e09020004040403040403010b2004200541016a220736029c0820052d0001220241ec004704402001418002722101200721050c050b2004200541026a220536029c0820014180067221010c030b2004200541016a220736029c0820052d0001220241e8004704402001418001722101200721050c040b2004200541026a220536029c08200141c0017221010c020b2004200541016a220536029c0820014180047221010c010b2004200541016a220536029c0820014180027221010b20052d000021020b0240027f024002400240024002400240024020024118744118752209419e7f6a220741164b044020094125470440200941c600460d07200941d800470d020c080b200341ff074b0d02200441106a20036a41253a00000c020b200741016b0e15040600050000060000000000060200000300060000060b200341ff074b0d00200441106a20036a20023a00000b200341016a21030c060b200441106a2003200028020041004110200841082001412172101621032004200541016a220536029c08200041046a21000c0b0b20002802002202417f6a21060340200641016a22062d00000d000b200620026b2206200820062008491b20062001418008712209410a761b210702402001410271220b0440200321060c010b200441106a20036a210c4100210103400240200120036a2106200120076a220d200a4f0d00200641ff074d04402001200c6a41203a00000b200141016a21010c010b0b200d41016a21070b200041046a2100034002402006210120022d00002203450d00200904402008450d012008417f6a21080b200141016a2106200241016a2102200141ff074b0d01200441106a20016a20033a00000c010b0b200b450440200121030c050b200441106a20016a2102410021060340200120066a2103200620076a200a4f0d05200341ff074d0440200220066a41203a00000b200641016a21060c000b000b410121060240200141027122070440200321020c010b200441106a20036a21084100210103400240200120036a2102200141016a2206200a4f0d00200241ff074d0440200120086a41203a00000b200621010c010b0b200141026a21060b200241ff074d0440200441106a20026a20002802003a00000b200241016a22032007450d021a0340200322012006200a4f0d031a200141016a2103200641016a2106200141ff074b0d00200441106a20016a41203a00000c000b000b200841062001418008711b220841037441a0086a2107200041076a417871220c2b0300211241002102034002402008410a490d002002411f4b0d00200441a0086a20026a41303a0000200741786a21072008417f6a2108200241016a21020c010b0b027f4400000000000000002012a12012201244000000000000000063220b1b22129944000000000000e0416304402012aa0c010b4180808080780b2100027f20122000b7a120072b03002214a2221344000000000000f041632013440000000000000000667104402013ab0c010b41000b2109024020132009b8a1221344000000000000e03f644101734504402014200941016a2209b8654101730d01200041016a2100410021090c010b201344000000000000e03f620d0020094520094101717220096a21090b4100210702402012440000c0ffffffdf41640d00024020080440200220086a41606a210f034002402002412046044041202102200f21080c010b200441a0086a20026a20092009410a6e220741766c6a4130723a0000200241016a21022008417f6a2108200941094b210d20072109200d0d010b0b034002402002411f4b22070d002008450d00200441a0086a20026a41303a0000200241016a21022008417f6a21080c010b0b20070d01200441a0086a20026a412e3a0000200241016a21020c010b20122000b7a1221244000000000000e03f64410173450440200041016a21000c010b20002000201244000000000000e03f61716a21000b03402002411f4d0440200441a0086a20026a20002000410a6d220841766c6a41306a3a0000200241016a2102200041096a210720082100200741124b0d010b0b20014103712108034002402002411f4b0d0020084101470d002002200a4f0d00200441a0086a20026a41303a0000200241016a21020c010b0b20014101712108200141027121000240200a2001410c71410047200b726b20022002200a461b2202411f4b0d000240200b410173450440200441a0086a20026a412d3a00000c010b20014104710440200441a0086a20026a412b3a00000c010b2001410871450d01200441a0086a20026a41203a00000b200241016a21020b024020002008720440200321070c010b200441106a20036a2108410021010340200120036a2107200120026a200a4f0d01200741ff074d0440200120086a41203a00000b200141016a21010c000b000b200620076a2103034020024504402000450d0203402003200a4f0d03200741ff074d0440200441106a20076a41203a00000b200741016a2107200341016a21030c000b000b200741ff074d0440200441106a20076a2002200e6a2d00003a00000b200741016a2107200341016a21032002417f6a21020c000b000b200c41086a21002004200541016a220536029c08200721030c080b41102107024002400240200241ff0171220641d8004622020d00200641f800460d00200641ef00460440410821070c020b200641e200460440410221070c020b2001416f712101410a21070b2001412072200120021b2101200641e400460d01200641e900460d010b200141737121010b2001417e7120012001418008711b2101200641e900474100200641e400471b4504402001418004710440200441106a2003200041076a417871220629030022102010423f8722117c2011852010423f88a72007ad2008200a200110172103200641086a21000c030b2001418002710440200441106a2003200028020022062006411f7522026a2002732006411f7620072008200a200110160c020b200441106a2003027f200141c00071044020002c00000c010b2000280200220641107441107520062001418001711b0b2206411f75220220066a2002732006411f7620072008200a200110160c010b2001418004710440200441106a2003200041076a417871220629030041002007ad2008200a200110172103200641086a21000c020b2001418002710440200441106a20032000280200410020072008200a200110160c010b200441106a2003027f200141c00071044020002d00000c010b2000280200220641ffff037120062001418001711b0b410020072008200a200110160b2103200041046a21000b2004200541016a220536029c080c050b2004200541016a220536029c08200141047221010c020b2004200541016a220536029c08200141027221010c010b2004200541016a220536029c08200141017221010c000b000b0b200441106a200341ff07200341ff07491b6a41003a0000200441106a2103024003402003410371044020032d0000450d02200341016a21030c010b0b2003417c6a21030340200341046a22032802002205417f73200541fffdfb776a7141808182847871450d000b0340200541ff0171450d01200341016a2d00002105200341016a21030c000b000b200441106a2003200441106a6b1001200441c0086a24000b990101047f230041106b220124002001200036020c2000047f418409200041086a220241107622004184092802006a2203360200418009418009280200220420026a41076a417871220236020002400240200341107420024d0440418409200341016a360200200041016a21000c010b2000450d010b200040000d0010000b20042001410c6a101241086a0541000b2100200141106a240020000ba50502087f017e230041406a220524001013100222001008220110030240200541206a20012000411c100a2204280208450440200441146a2802002100200428021021010c010b200541386a2004100b20042005280238200528023c100c36020c200541086a2004100b410021002004027f410020052802082202450d001a4100200528020c2206200428020c2203490d001a200620032003417f461b210020020b2201360210200441146a2000360200200441003602080b200541086a200120004114100a2200100d024002402000280204450d002000100d0240200028020022022c0000220141004e044020010d010c020b200141807f460d00200141ff0171220341b7014d0440200028020441014d04401000200028020021020b20022d00010d010c020b200341bf014b0d012000280204200141ff017141ca7e6a22014d04401000200028020021020b200120026a2d0000450d010b2000280204450d0020022d000041c001490d010b10000b2000100e2206200028020422024b04401000200028020421020b20002802002107024002400240200204404100210320072c00002200417f4a0d01027f200041ff0171220341bf014d04404100200041ff017141b801490d011a200341c97e6a0c010b4100200041ff017141f801490d001a200341897e6a0b41016a21030c010b4101210320070d00410021000c010b41002100200320066a20024b0d0020022006490d004100210120022003490d01200320076a2101200220036b20062006417f461b22004109490d0110000c010b410021010b0340200004402000417f6a210020013100002008420886842108200141016a21010c010b0b024002402008500d00418908100f20085104402004410210100c020b418e08100f2008520d002004410310100c010b10000b1011200541406b24000b730020004200370210200042ffffffff0f370208200020023602042000200136020002402003410871450d002000101920024f0d002003410471044010000c010b200042003702000b02402003411071450d002000101920024d0d0020034104710440100020000f0b200042003702000b20000b7201047f2001100e220220012802044b044010000b2001101a21032000027f0240200128020022054504400c010b200220036a200128020422014b0d0020012003490d00410020012002490d011a200320056a2104200120036b20022002417f461b0c010b41000b360204200020043602000b2701017f230041206b22022400200241086a200020014114100a10192100200241206a240020000b4101017f200028020445044010000b0240200028020022012d0000418101470d00200028020441014d047f100020002802000520010b2c00014100480d0010000b0bff0201037f200028020445044041000f0b2000100d41012102024020002802002c00002201417f4a0d00200141ff0171220341b7014d0440200341807f6a0f0b02400240200141ff0171220141bf014d0440024020002802042201200341c97e6a22024d047f100020002802040520010b4102490d0020002802002d00010d0010000b200241054f044010000b20002802002d000145044010000b4100210241b7012101034020012003460440200241384f0d030c0405200028020020016a41ca7e6a2d00002002410874722102200141016a21010c010b000b000b200141f7014d0440200341c07e6a0f0b024020002802042201200341897e6a22024d047f100020002802040520010b4102490d0020002802002d00010d0010000b200241054f044010000b20002802002d000145044010000b4100210241f701210103402001200346044020024138490d0305200028020020016a418a7e6a2d00002002410874722102200141016a21010c010b0b0b200241ff7d490d010b10000b20020b3901027e42a5c688a1c89ca7f94b210103402000300000220250450440200041016a2100200142b383808080207e20028521010c010b0b20010b920101047f230041106b22022400024002402000280204450d0020002802002d000041c001490d00200241086a2000100b41012104200228020c2100034020000440200241002002280208220320032000100c22056a20034520002005497222031b3602084100200020056b20031b21002004417f6a21040c010b0b2004450d010b10000b20022001110000200241106a24000b880101037f41f008410136020041f4082802002100034020000440034041f80841f8082802002201417f6a2202360200200141014845044041f0084100360200200020024102746a22004184016a280200200041046a28020011000041f008410136020041f40828020021000c010b0b41f808412036020041f408200028020022003602000c010b0b0bf50801077f03400240200020036a2104200120036a210220034104460d002002410371450d00200420022d00003a0000200341016a21030c010b0b410420036b21050240200441037122064504400340200541104f0440200020036a2202200120036a2204290200370200200241086a200441086a290200370200200341106a2103200541706a21050c010b0b027f2005410871450440200120036a2102200020036a0c010b200020036a2204200120036a2203290200370200200341086a2102200441086a0b21032005410471044020032002280200360200200341046a2103200241046a21020b20054102710440200320022f00003b0000200341026a2103200241026a21020b2005410171450d01200320022d00003a000020000f0b024020054120490d002006417f6a220641024b0d00024002400240024002400240200641016b0e020102000b2004200120036a220228020022063a0000200441016a200241016a2f00003b0000200041036a2108410120036b2105034020054111490d03200320086a2202200120036a220441046a2802002207410874200641187672360200200241046a200441086a2802002206410874200741187672360200200241086a2004410c6a28020022074108742006411876723602002002410c6a200441106a2802002206410874200741187672360200200341106a2103200541706a21050c000b000b2004200120036a220228020022063a0000200441016a200241016a2d00003a0000200041026a2108410220036b2105034020054112490d03200320086a2202200120036a220441046a2802002207411074200641107672360200200241046a200441086a2802002206411074200741107672360200200241086a2004410c6a28020022074110742006411076723602002002410c6a200441106a2802002206411074200741107672360200200341106a2103200541706a21050c000b000b2004200120036a28020022063a0000200041016a21082003417f7341046a2105034020054113490d03200320086a2202200120036a220441046a2802002207411874200641087672360200200241046a200441086a2802002206411874200741087672360200200241086a2004410c6a28020022074118742006410876723602002002410c6a200441106a2802002206411874200741087672360200200341106a2103200541706a21050c000b000b200120036a41036a2102200020036a41036a21040c020b200120036a41026a2102200020036a41026a21040c010b200120036a41016a2102200020036a41016a21040b20054110710440200420022d00003a00002004200228000136000120042002290005370005200420022f000d3b000d200420022d000f3a000f200441106a2104200241106a21020b2005410871044020042002290000370000200441086a2104200241086a21020b2005410471044020042002280000360000200441046a2104200241046a21020b20054102710440200420022f00003b0000200441026a2104200241026a21020b2005410171450d00200420022d00003a00000b20000b3501017f230041106b22004190890436020c41fc08200028020c41076a417871220036020041800920003602004184093f003602000b4501037f20002802002101034020012d000041506a41ff017141094b4504402000200141016a220336020020012c00002002410a6c6a41506a2102200321010c010b0b20020b140020022003490440200120026a20003a00000b0ba70101067f230041206b220924000240200245410020072007416f7120021b220a418008711b0d00200a41207141e1007341f6016a210b410021070340200720096a2002200220046e220c20046c6b22084130200b200841187441808080d000481b6a3a0000200741016a21082007411e4b0d01200220044f210d20082107200c2102200d0d000b0b20002001200920082003200420052006200a10182102200941206a240020020bb60102057f017e230041206b2209240020072007416f71200242005222081b210a0240200845044041002108200a418008710d010b200a41207141e1007341f6016a210b410021070340200720096a4130200b20022002200480220d20047e7da7220841187441808080d000481b20086a3a0000200741016a21082007411e4b0d01200220045a210c20082107200d2102200c0d000b0b200020012009200820032004a720052006200a10182107200941206a240020070bec0401037f2008410271210b2003210902400340200b0d010240200920064f0d002009411f4b0d00200220096a41303a0000200941016a21090c010b0b200921030b2008410371410147210a2003210902400340200a0d010240200920074f0d002009411f4b0d00200220096a41303a0000200941016a21090c010b0b200921030b2008410171210a0240024002402008411071044002402008418008710d002003450d002003200647410020032007471b0d002003417e6a2003417f6a220920091b200920054110461b21030b0240200541104604400240200841207122090d002003411f4b0d00200220036a41f8003a0000200341016a21030c020b2009450d012003411f4b0d01200220036a41d8003a0000200341016a21030c010b20054102470d002003411f4b0d00200220036a41e2003a0000200341016a21030b2003411f4b0d01200220036a41303a0000200341016a21030c010b20030d00410021030c010b20072008410c714100472004726b200320032007461b2203411f4b0d010b20040440200220036a412d3a0000200341016a21030c010b20084104710440200220036a412b3a0000200341016a21030c010b2008410871450d00200220036a41203a0000200341016a21030b200121090240200a200b720d002003210a0340200a20074f0d014120200020094180081015200a41016a210a200941016a21090c000b000b2002417f6a210a037f2003047f2003200a6a2c00002000200941800810152003417f6a2103200941016a21090c01050240200b450d00410020016b21030340200320096a20074f0d014120200020094180081015200941016a21090c000b000b20090b0b0b2e01017f200028020445044041000f0b4101210120002802002c0000417f4c047f2000101a2000100e6a0520010b0b5b00027f027f41002000280204450d001a410020002802002c0000417f4a0d011a20002802002d0000220041bf014d04404100200041b801490d011a200041c97e6a0c010b4100200041f801490d001a200041897e6a0b41016a0b0b0b7102004180080b1a69203d202564090a00696e6974006d656d6f72795f6c696d69740041a6080b4af03f000000000000244000000000000059400000000000408f40000000000088c34000000000006af8400000000080842e4100000000d01263410000000084d797410000000065cdcd41";

    public static String BINARY = BINARY_0;

    public static final String FUNC_MEMORY_LIMIT = "memory_limit";

    protected OOMException(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected OOMException(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<OOMException> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(OOMException.class, web3j, credentials, contractGasProvider, encodedConstructor, chainId);
    }

    public static RemoteCall<OOMException> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(OOMException.class, web3j, transactionManager, contractGasProvider, encodedConstructor, chainId);
    }

    public static RemoteCall<OOMException> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(OOMException.class, web3j, credentials, contractGasProvider, encodedConstructor, initialVonValue, chainId);
    }

    public static RemoteCall<OOMException> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(OOMException.class, web3j, transactionManager, contractGasProvider, encodedConstructor, initialVonValue, chainId);
    }

    public RemoteCall<TransactionReceipt> memory_limit() {
        final WasmFunction function = new WasmFunction(FUNC_MEMORY_LIMIT, Arrays.asList(), Void.class);
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> memory_limit(BigInteger vonValue) {
        final WasmFunction function = new WasmFunction(FUNC_MEMORY_LIMIT, Arrays.asList(), Void.class);
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static OOMException load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new OOMException(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static OOMException load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new OOMException(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
